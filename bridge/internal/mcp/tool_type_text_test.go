package mcp

import "testing"

func TestTypeTextTool(t *testing.T) {
	result, mock := callTool(t, "tv_type_text", map[string]any{"text": "stranger things"}, everyScope)
	if isError(result) || mock.Calls[0] != "text" {
		t.Fatalf("unexpected calls: %v", mock.Calls)
	}
}

func TestTypeTextToolRejectsShellCharacters(t *testing.T) {
	for _, text := range []string{"a; reboot", "$(id)", "`id`", "a|b", "a&b", "quote'", "new\nline", ""} {
		result, mock := callTool(t, "tv_type_text", map[string]any{"text": text}, everyScope)
		if !isError(result) || len(mock.Calls) != 0 {
			t.Fatalf("text %q must be rejected", text)
		}
	}
}
