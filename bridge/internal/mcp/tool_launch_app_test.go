package mcp

import "testing"

func TestLaunchAppTool(t *testing.T) {
	result, mock := callTool(t, "tv_launch_app", map[string]any{"package": "com.netflix.ninja"}, everyScope)
	if isError(result) || mock.Calls[0] != "launch com.netflix.ninja" {
		t.Fatalf("unexpected calls: %v", mock.Calls)
	}
}

func TestLaunchAppToolRejectsInjectionAndUnknownApps(t *testing.T) {
	for _, name := range []string{"com.netflix.ninja;reboot", "$(id)", "../../etc", "com.not.installed"} {
		result, mock := callTool(t, "tv_launch_app", map[string]any{"package": name}, everyScope)
		if !isError(result) || len(mock.Calls) != 0 {
			t.Fatalf("package %q must be rejected", name)
		}
	}
}
