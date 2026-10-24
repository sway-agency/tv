package mcp

import (
	"context"
	"encoding/json"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/tv"
)

func powerTool() Tool {
	return Tool{
		Name:        "tv_power",
		Title:       "Turn the screen on or off",
		Description: "Wakes the TV screen or puts it to standby. Some TVs drop off the network in deep standby and cannot be woken remotely.",
		Scope:       auth.ScopeControl,
		InputSchema: objectSchema(map[string]any{
			"on": map[string]any{"type": "boolean"},
		}, "on"),
		Run: func(ctx context.Context, controller tv.Controller, arguments json.RawMessage) (ToolResult, error) {
			var input struct {
				On *bool `json:"on"`
			}
			err := decodeArguments(arguments, &input)
			if err != nil {
				return ToolResult{}, err
			}
			if input.On == nil {
				return ToolResult{}, errMissing("on")
			}
			err = controller.SetPower(ctx, *input.On)
			if err != nil {
				return ToolResult{}, err
			}
			if *input.On {
				return ToolResult{Text: "Screen on."}, nil
			}
			return ToolResult{Text: "Screen in standby."}, nil
		},
	}
}
