package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/config"
	"github.com/sway-agency/tv/bridge/internal/mcp"
)

const (
	maxFailures     = 10
	failureWindow   = 10 * time.Minute
	lockoutDuration = 30 * time.Minute
)

func New(settings config.Config, dispatcher *mcp.Dispatcher, store *auth.Store, logger *slog.Logger) (*http.Server, error) {
	tlsSettings, err := tlsConfig(settings)
	if err != nil {
		return nil, err
	}
	handlers := &handler{
		dispatcher: dispatcher,
		store:      store,
		limiter:    auth.NewLimiter(settings.RateLimitPerSecond, settings.RateLimitBurst),
		lockout:    auth.NewLockout(maxFailures, failureWindow, lockoutDuration),
		logger:     logger,
	}
	routes := http.NewServeMux()
	routes.HandleFunc("/mcp", handlers.serveMCP)
	routes.HandleFunc("POST /v1/tools/{name}", handlers.serveTool)
	return &http.Server{
		Addr:              settings.Listen,
		Handler:           guard(settings, routes),
		TLSConfig:         tlsSettings,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      45 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    16 << 10,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
	}, nil
}

func tlsConfig(settings config.Config) (*tls.Config, error) {
	pair, err := tls.LoadX509KeyPair(settings.CertFile, settings.KeyFile)
	if err != nil {
		return nil, err
	}
	result := &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{pair},
		NextProtos:   []string{"h2", "http/1.1"},
	}
	if !settings.RequireClientCert {
		return result, nil
	}
	authority, err := os.ReadFile(settings.ClientCAFile)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(authority) {
		return nil, errors.New("client_ca_file holds no certificate")
	}
	result.ClientCAs = pool
	result.ClientAuth = tls.RequireAndVerifyClientCert
	return result, nil
}
