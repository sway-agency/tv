package server

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/sway-agency/tv/bridge/internal/auth"
)

const modernList = `{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{"_meta":{"io.modelcontextprotocol/protocolVersion":"2026-07-28","io.modelcontextprotocol/clientCapabilities":{}}}}`

var modernHeaders = map[string]string{"MCP-Protocol-Version": "2026-07-28", "Mcp-Method": "tools/list"}

func TestModernToolsListOverTLS(t *testing.T) {
	bridge := startBridge(t, auth.ScopeObserve)
	response := bridge.post(t, "/mcp", modernList, modernHeaders)
	body, _ := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK || !strings.Contains(string(body), "tv_status") || strings.Contains(string(body), "tv_press_key") {
		t.Fatalf("unexpected response %d %s", response.StatusCode, body)
	}
	if response.TLS.Version != tls.VersionTLS13 || response.Header.Get("Strict-Transport-Security") == "" {
		t.Fatal("TLS 1.3 and security headers expected")
	}
}

func TestTLS12IsRefused(t *testing.T) {
	bridge := startBridge(t, auth.ScopeObserve)
	transport := bridge.client.Transport.(*http.Transport)
	transport.TLSClientConfig.MaxVersion = tls.VersionTLS12
	_, err := bridge.client.Post(bridge.url+"/mcp", "application/json", strings.NewReader(modernList))
	if err == nil {
		t.Fatal("TLS 1.2 handshake must fail")
	}
}

func TestBadTokensLockTheAddressOut(t *testing.T) {
	bridge := startBridge(t, auth.ScopeObserve)
	valid := bridge.token
	bridge.token = "acb_not-a-real-token"
	for range maxFailures {
		response := bridge.post(t, "/mcp", modernList, modernHeaders)
		if response.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", response.StatusCode)
		}
	}
	bridge.token = valid
	response := bridge.post(t, "/mcp", modernList, modernHeaders)
	if response.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("locked address must be refused even with a valid token, got %d", response.StatusCode)
	}
}

func TestGuardsRejectForeignHostOriginMethodAndContentType(t *testing.T) {
	bridge := startBridge(t, auth.ScopeObserve)
	cases := []struct {
		headers map[string]string
		status  int
	}{
		{map[string]string{"Origin": "https://evil.example"}, http.StatusForbidden},
		{map[string]string{"Host": "tv.attacker.example"}, http.StatusMisdirectedRequest},
		{map[string]string{"Content-Type": "text/plain"}, http.StatusBadRequest},
	}
	for _, testCase := range cases {
		request, _ := http.NewRequest(http.MethodPost, bridge.url+"/mcp", strings.NewReader(modernList))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+bridge.token)
		for name, value := range testCase.headers {
			request.Header.Set(name, value)
		}
		request.Host = request.Header.Get("Host")
		response, err := bridge.client.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
		if response.StatusCode != testCase.status {
			t.Fatalf("%v: expected %d, got %d", testCase.headers, testCase.status, response.StatusCode)
		}
	}
	response, err := bridge.client.Get(bridge.url + "/mcp")
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("GET must return 405, got %d", response.StatusCode)
	}
}

func TestRESTToolCall(t *testing.T) {
	bridge := startBridge(t, auth.ScopeObserve, auth.ScopeControl)
	response := bridge.post(t, "/v1/tools/tv_launch_app", `{"package":"com.netflix.ninja"}`, nil)
	if response.StatusCode != http.StatusOK || len(bridge.mock.Calls) == 0 {
		t.Fatalf("REST call failed: %d %v", response.StatusCode, bridge.mock.Calls)
	}
}
