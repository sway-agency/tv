package mcp

import "testing"

func TestStatusTool(t *testing.T) {
	result, mock := callTool(t, "tv_status", nil, everyScope)
	if isError(result) || mock.Calls[0] != "status" || result["structuredContent"] == nil {
		t.Fatalf("unexpected status result: %+v", result)
	}
}

func TestStatusToolRejectsUnknownArguments(t *testing.T) {
	result, mock := callTool(t, "tv_status", map[string]any{"debug": true}, everyScope)
	if !isError(result) || len(mock.Calls) != 0 {
		t.Fatalf("unknown arguments must be rejected")
	}
}
