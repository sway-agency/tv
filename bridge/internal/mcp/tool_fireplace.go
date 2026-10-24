package mcp

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/tv"
)

func fireplaceTool() Tool {
	return Tool{
		Name:        "tv_fireplace",
		Title:       "Start the fireplace",
		Description: "Shows the crackling fireplace full screen. Press back to leave it.",
		Scope:       auth.ScopeControl,
		InputSchema: objectSchema(map[string]any{}),
		Run: func(ctx context.Context, controller tv.Controller, arguments json.RawMessage) (ToolResult, error) {
			err := decodeArguments(arguments, &struct{}{})
			if err != nil {
				return ToolResult{}, err
			}
			err = controller.StartFireplace(ctx)
			if err != nil {
				return ToolResult{}, err
			}
			return ToolResult{Text: "Fireplace started."}, nil
		},
	}
}

func errMissing(field string) error {
	return errors.New("missing required argument: " + field)
}
