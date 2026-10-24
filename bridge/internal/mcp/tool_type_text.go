package mcp

import (
	"context"
	"encoding/json"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/tv"
)

func typeTextTool() Tool {
	return Tool{
		Name:        "tv_type_text",
		Title:       "Type text",
		Description: "Types text into the focused field, for example a search box. Letters, digits and spaces only.",
		Scope:       auth.ScopeText,
		InputSchema: objectSchema(map[string]any{
			"text": map[string]any{"type": "string", "minLength": 1, "maxLength": tv.MaxTextLength, "pattern": "^[a-zA-Z0-9 ]+$"},
		}, "text"),
		Run: func(ctx context.Context, controller tv.Controller, arguments json.RawMessage) (ToolResult, error) {
			var input struct {
				Text string `json:"text"`
			}
			err := decodeArguments(arguments, &input)
			if err != nil {
				return ToolResult{}, err
			}
			err = tv.ValidateText(input.Text)
			if err != nil {
				return ToolResult{}, err
			}
			err = controller.TypeText(ctx, input.Text)
			if err != nil {
				return ToolResult{}, err
			}
			return ToolResult{Text: "Text typed."}, nil
		},
	}
}
