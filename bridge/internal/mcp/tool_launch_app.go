package mcp

import (
	"context"
	"encoding/json"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/tv"
)

func launchAppTool() Tool {
	return Tool{
		Name:        "tv_launch_app",
		Title:       "Launch a TV app",
		Description: "Opens an installed TV app by package name. Call tv_list_apps to get the names.",
		Scope:       auth.ScopeControl,
		InputSchema: objectSchema(map[string]any{
			"package": map[string]any{"type": "string", "maxLength": 255, "pattern": `^[a-zA-Z][a-zA-Z0-9_]*(\.[a-zA-Z][a-zA-Z0-9_]*)+$`},
		}, "package"),
		Run: func(ctx context.Context, controller tv.Controller, arguments json.RawMessage) (ToolResult, error) {
			var input struct {
				Package string `json:"package"`
			}
			err := decodeArguments(arguments, &input)
			if err != nil {
				return ToolResult{}, err
			}
			err = tv.ValidatePackage(input.Package)
			if err != nil {
				return ToolResult{}, err
			}
			err = controller.LaunchApp(ctx, input.Package)
			if err != nil {
				return ToolResult{}, err
			}
			return ToolResult{Text: "Launched " + input.Package + "."}, nil
		},
	}
}
