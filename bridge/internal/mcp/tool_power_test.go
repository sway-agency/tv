package mcp

import "testing"

func TestPowerTool(t *testing.T) {
	result, mock := callTool(t, "tv_power", map[string]any{"on": false}, everyScope)
	if isError(result) || mock.Calls[0] != "power false" {
		t.Fatalf("unexpected calls: %v", mock.Calls)
	}
}

func TestPowerToolRequiresOn(t *testing.T) {
	result, _ := callTool(t, "tv_power", map[string]any{}, everyScope)
	if !isError(result) {
		t.Fatalf("missing on must be rejected")
	}
}
