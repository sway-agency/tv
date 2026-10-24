package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/tv"
)

func statusTool() Tool {
	return Tool{
		Name:        "tv_status",
		Title:       "TV status",
		Description: "Tells whether the TV screen is on and which app is in the foreground.",
		Scope:       auth.ScopeObserve,
		ReadOnly:    true,
		InputSchema: objectSchema(map[string]any{}),
		Run: func(ctx context.Context, controller tv.Controller, arguments json.RawMessage) (ToolResult, error) {
			err := decodeArguments(arguments, &struct{}{})
			if err != nil {
				return ToolResult{}, err
			}
			status, err := controller.Status(ctx)
			if err != nil {
				return ToolResult{}, err
			}
			text := fmt.Sprintf("Screen on: %t. Foreground app: %s.", status.Awake, status.ForegroundPackage)
			return ToolResult{Text: text, Structured: status}, nil
		},
	}
}
