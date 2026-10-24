package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func handle(request Request, envelope Envelope) Outcome {
	dispatcher, _ := newTestDispatcher()
	return dispatcher.Handle(context.Background(), request, envelope, everyScope)
}

func TestModernDiscover(t *testing.T) {
	outcome := handle(Request{JSONRPC: "2.0", ID: json.RawMessage(`"a"`), Method: "server/discover", Params: modernParams(nil)},
		Envelope{HTTP: true, ProtocolVersion: ModernVersion, Method: "server/discover"})
	result, _ := outcome.Response.Result.(map[string]any)
	if outcome.Status != http.StatusOK || result["resultType"] != "complete" {
		t.Fatalf("unexpected discover outcome: %+v", outcome.Response)
	}
}

func TestModernRejectsHeaderMismatch(t *testing.T) {
	outcome := handle(Request{JSONRPC: "2.0", ID: json.RawMessage("1"), Method: "tools/list", Params: modernParams(nil)},
		Envelope{HTTP: true, ProtocolVersion: ModernVersion, Method: "tools/call"})
	if outcome.Status != http.StatusBadRequest || outcome.Response.Error.Code != codeHeaderMismatch {
		t.Fatalf("expected header mismatch, got %+v", outcome.Response.Error)
	}
}

func TestModernRejectsMissingCapabilities(t *testing.T) {
	params := json.RawMessage(`{"_meta":{"io.modelcontextprotocol/protocolVersion":"2026-07-28"}}`)
	outcome := handle(Request{JSONRPC: "2.0", ID: json.RawMessage("1"), Method: "tools/list", Params: params}, Envelope{})
	if outcome.Status != http.StatusBadRequest || outcome.Response.Error.Code != codeInvalidParams {
		t.Fatalf("expected invalid params, got %+v", outcome.Response.Error)
	}
}

func TestUnsupportedVersion(t *testing.T) {
	params := json.RawMessage(`{"_meta":{"io.modelcontextprotocol/protocolVersion":"1900-01-01","io.modelcontextprotocol/clientCapabilities":{}}}`)
	outcome := handle(Request{JSONRPC: "2.0", ID: json.RawMessage("1"), Method: "tools/list", Params: params}, Envelope{})
	if outcome.Response.Error.Code != codeUnsupportedVersion {
		t.Fatalf("expected unsupported version, got %+v", outcome.Response.Error)
	}
}

func TestModernUnknownMethodIsNotFound(t *testing.T) {
	outcome := handle(Request{JSONRPC: "2.0", ID: json.RawMessage("1"), Method: "resources/list", Params: modernParams(nil)}, Envelope{})
	if outcome.Status != http.StatusNotFound || outcome.Response.Error.Code != codeMethodNotFound {
		t.Fatalf("expected method not found, got %+v", outcome)
	}
}
