package mcp

import "testing"

func TestScreenshotTool(t *testing.T) {
	result, mock := callTool(t, "tv_screenshot", nil, everyScope)
	content, _ := result["content"].([]any)
	if isError(result) || mock.Calls[0] != "screenshot" || len(content) != 1 {
		t.Fatalf("unexpected screenshot result: %+v", result)
	}
}
