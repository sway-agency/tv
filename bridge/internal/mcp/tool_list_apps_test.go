package mcp

import "testing"

func TestListAppsTool(t *testing.T) {
	result, _ := callTool(t, "tv_list_apps", map[string]any{}, everyScope)
	if isError(result) {
		t.Fatalf("unexpected error: %+v", result)
	}
}
