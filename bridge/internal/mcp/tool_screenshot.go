package mcp

import (
	"context"
	"encoding/json"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/tv"
)

func screenshotTool() Tool {
	return Tool{
		Name:        "tv_screenshot",
		Title:       "Screenshot",
		Description: "Captures what the TV screen shows, as a small JPEG. Only granted to tokens with the screen scope.",
		Scope:       auth.ScopeScreen,
		ReadOnly:    true,
		InputSchema: objectSchema(map[string]any{}),
		Run: func(ctx context.Context, controller tv.Controller, arguments json.RawMessage) (ToolResult, error) {
			err := decodeArguments(arguments, &struct{}{})
			if err != nil {
				return ToolResult{}, err
			}
			image, err := controller.Screenshot(ctx)
			if err != nil {
				return ToolResult{}, err
			}
			return ToolResult{ImageJPEG: image}, nil
		},
	}
}
