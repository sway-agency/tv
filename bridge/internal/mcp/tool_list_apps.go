package mcp

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/tv"
)

func listAppsTool() Tool {
	return Tool{
		Name:        "tv_list_apps",
		Title:       "List TV apps",
		Description: "Lists the package names of the apps that can be launched on the TV.",
		Scope:       auth.ScopeObserve,
		ReadOnly:    true,
		InputSchema: objectSchema(map[string]any{}),
		Run: func(ctx context.Context, controller tv.Controller, arguments json.RawMessage) (ToolResult, error) {
			err := decodeArguments(arguments, &struct{}{})
			if err != nil {
				return ToolResult{}, err
			}
			apps, err := controller.Apps(ctx)
			if err != nil {
				return ToolResult{}, err
			}
			names := make([]string, 0, len(apps))
			for _, app := range apps {
				names = append(names, app.Package)
			}
			return ToolResult{Text: strings.Join(names, "\n"), Structured: map[string]any{"apps": apps}}, nil
		},
	}
}
