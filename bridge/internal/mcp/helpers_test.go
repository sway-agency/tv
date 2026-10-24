package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/tv"
)

var everyScope = Caller{Name: "test", Scopes: auth.AllScopes}

func newTestDispatcher() (*Dispatcher, *tv.Mock) {
	mock := tv.NewMock("com.netflix.ninja", "net.oqee.androidtv.store")
	return NewDispatcher(mock, "test"), mock
}

func modernParams(extra map[string]any) json.RawMessage {
	params := map[string]any{"_meta": map[string]any{
		metaVersion:      ModernVersion,
		metaCapabilities: map[string]any{},
	}}
	for key, value := range extra {
		params[key] = value
	}
	encoded, _ := json.Marshal(params)
	return encoded
}

func callTool(t *testing.T, name string, arguments map[string]any, caller Caller) (map[string]any, *tv.Mock) {
	t.Helper()
	dispatcher, mock := newTestDispatcher()
	params, _ := json.Marshal(map[string]any{"name": name, "arguments": arguments})
	outcome := dispatcher.Handle(context.Background(), Request{JSONRPC: "2.0", ID: json.RawMessage("1"), Method: "tools/call", Params: params}, Envelope{}, caller)
	if outcome.Response == nil || outcome.Response.Error != nil {
		t.Fatalf("tools/call %s failed: %+v", name, outcome.Response)
	}
	result, _ := outcome.Response.Result.(map[string]any)
	return result, mock
}

func isError(result map[string]any) bool {
	failed, _ := result["isError"].(bool)
	return failed
}
