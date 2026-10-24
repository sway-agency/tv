package server

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/config"
	"github.com/sway-agency/tv/bridge/internal/mcp"
	"github.com/sway-agency/tv/bridge/internal/pki"
	"github.com/sway-agency/tv/bridge/internal/tv"
)

type testBridge struct {
	url    string
	token  string
	client *http.Client
	mock   *tv.Mock
}

func startBridge(t *testing.T, scopes ...auth.Scope) testBridge {
	t.Helper()
	dir := t.TempDir()
	authority, err := pki.NewAuthority(dir)
	if err != nil {
		t.Fatal(err)
	}
	err = authority.IssueServer(dir, []string{"127.0.0.1"})
	if err != nil {
		t.Fatal(err)
	}
	plain, token, err := auth.NewToken("agent", scopes)
	if err != nil {
		t.Fatal(err)
	}
	settings := config.Config{
		Listen: "127.0.0.1:0", Hosts: []string{"127.0.0.1"}, AllowedOrigins: []string{},
		CertFile: filepath.Join(dir, "server.pem"), KeyFile: filepath.Join(dir, "server-key.pem"),
		RateLimitPerSecond: 50, RateLimitBurst: 50, Tokens: []auth.Token{token},
	}
	mock := tv.NewMock("com.netflix.ninja")
	httpServer, err := New(settings, mcp.NewDispatcher(mock, "test"), slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go func() { _ = httpServer.ServeTLS(listener, "", "") }()
	t.Cleanup(func() { _ = httpServer.Close() })
	return testBridge{url: "https://" + listener.Addr().String(), token: plain, client: trustingClient(t, dir, tls.VersionTLS13), mock: mock}
}

func trustingClient(t *testing.T, dir string, maxVersion uint16) *http.Client {
	t.Helper()
	authority, err := os.ReadFile(filepath.Join(dir, "ca.pem"))
	if err != nil {
		t.Fatal(err)
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(authority)
	return &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: pool, MaxVersion: maxVersion}}}
}

func (b testBridge) post(t *testing.T, path string, body string, headers map[string]string) *http.Response {
	t.Helper()
	request, err := http.NewRequest(http.MethodPost, b.url+path, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+b.token)
	for name, value := range headers {
		request.Header.Set(name, value)
	}
	response, err := b.client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { response.Body.Close() })
	return response
}
