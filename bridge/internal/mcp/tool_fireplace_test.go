package mcp

import "testing"

func TestFireplaceTool(t *testing.T) {
	result, mock := callTool(t, "tv_fireplace", nil, everyScope)
	if isError(result) || mock.Calls[0] != "fireplace" {
		t.Fatalf("unexpected calls: %v", mock.Calls)
	}
}
