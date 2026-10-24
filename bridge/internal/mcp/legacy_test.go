package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/sway-agency/tv/bridge/internal/auth"
)

func TestLegacyInitializeEchoesKnownVersion(t *testing.T) {
	outcome := handle(Request{JSONRPC: "2.0", ID: json.RawMessage("0"), Method: "initialize",
		Params: json.RawMessage(`{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"c","version":"1"}}`)}, Envelope{HTTP: true})
	result, _ := outcome.Response.Result.(map[string]any)
	if result["protocolVersion"] != "2025-06-18" {
		t.Fatalf("expected echoed version, got %v", result["protocolVersion"])
	}
}

func TestLegacyNotificationIsAccepted(t *testing.T) {
	outcome := handle(Request{JSONRPC: "2.0", Method: "notifications/initialized"}, Envelope{HTTP: true})
	if outcome.Status != http.StatusAccepted || outcome.Response != nil {
		t.Fatalf("expected 202 without body, got %+v", outcome)
	}
}

func TestLegacyToolsListFiltersByScope(t *testing.T) {
	dispatcher, _ := newTestDispatcher()
	caller := Caller{Name: "observer", Scopes: []auth.Scope{auth.ScopeObserve}}
	outcome := dispatcher.Handle(context.Background(), Request{JSONRPC: "2.0", ID: json.RawMessage("2"), Method: "tools/list"},
		Envelope{HTTP: true, ProtocolVersion: "2025-11-25"}, caller)
	result, _ := outcome.Response.Result.(map[string]any)
	tools, _ := result["tools"].([]map[string]any)
	if len(tools) != 2 {
		t.Fatalf("observe scope should list 2 tools, got %d", len(tools))
	}
}

func TestNullIDIsRejected(t *testing.T) {
	outcome := handle(Request{JSONRPC: "2.0", ID: json.RawMessage("null"), Method: "tools/list"}, Envelope{})
	if outcome.Status != http.StatusBadRequest {
		t.Fatalf("expected 400 for null id, got %d", outcome.Status)
	}
}

func TestToolOutsideScopeIsUnknown(t *testing.T) {
	dispatcher, mock := newTestDispatcher()
	params := json.RawMessage(`{"name":"tv_screenshot","arguments":{}}`)
	caller := Caller{Name: "limited", Scopes: []auth.Scope{auth.ScopeObserve, auth.ScopeControl}}
	outcome := dispatcher.Handle(context.Background(), Request{JSONRPC: "2.0", ID: json.RawMessage("3"), Method: "tools/call", Params: params}, Envelope{}, caller)
	if outcome.Response.Error == nil || len(mock.Calls) != 0 {
		t.Fatalf("screenshot must be refused without the screen scope")
	}
}
