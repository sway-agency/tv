package mcp

import "testing"

func TestPressKeyTool(t *testing.T) {
	result, mock := callTool(t, "tv_press_key", map[string]any{"key": "volume_up", "repeat": 3}, everyScope)
	if isError(result) || mock.Calls[0] != "key volume_up x3" {
		t.Fatalf("unexpected calls: %v", mock.Calls)
	}
}

func TestPressKeyToolRejectsUnknownKeyAndRepeat(t *testing.T) {
	for _, arguments := range []map[string]any{
		{"key": "power; reboot"},
		{"key": "up", "repeat": 500},
		{"key": "up", "repeat": "2"},
	} {
		result, mock := callTool(t, "tv_press_key", arguments, everyScope)
		if !isError(result) || len(mock.Calls) != 0 {
			t.Fatalf("arguments %v must be rejected", arguments)
		}
	}
}
