package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/tv"
)

func pressKeyTool() Tool {
	return Tool{
		Name:        "tv_press_key",
		Title:       "Press a remote key",
		Description: "Presses a remote control key, optionally several times (navigation, playback, volume).",
		Scope:       auth.ScopeControl,
		InputSchema: objectSchema(map[string]any{
			"key":    map[string]any{"type": "string", "enum": tv.KeyNames()},
			"repeat": map[string]any{"type": "integer", "minimum": 1, "maximum": tv.MaxKeyRepeat, "default": 1},
		}, "key"),
		Run: func(ctx context.Context, controller tv.Controller, arguments json.RawMessage) (ToolResult, error) {
			var input struct {
				Key    string `json:"key"`
				Repeat int    `json:"repeat"`
			}
			err := decodeArguments(arguments, &input)
			if err != nil {
				return ToolResult{}, err
			}
			key, err := tv.ParseKey(input.Key)
			if err != nil {
				return ToolResult{}, err
			}
			repeat := max(input.Repeat, 1)
			if repeat > tv.MaxKeyRepeat {
				return ToolResult{}, fmt.Errorf("repeat must be at most %d", tv.MaxKeyRepeat)
			}
			err = controller.PressKey(ctx, key, repeat)
			if err != nil {
				return ToolResult{}, err
			}
			return ToolResult{Text: fmt.Sprintf("Pressed %s %d time(s).", key.Name, repeat)}, nil
		},
	}
}
