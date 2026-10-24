package mcp

import (
	"context"
	"encoding/json"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/tv"
)

type ToolResult struct {
	Text       string
	Structured any
	ImageJPEG  []byte
}

type Tool struct {
	Name        string
	Title       string
	Description string
	Scope       auth.Scope
	ReadOnly    bool
	InputSchema map[string]any
	Run         func(ctx context.Context, controller tv.Controller, arguments json.RawMessage) (ToolResult, error)
}

func DefaultTools() []Tool {
	return []Tool{
		statusTool(), listAppsTool(), pressKeyTool(), launchAppTool(),
		powerTool(), fireplaceTool(), typeTextTool(), screenshotTool(),
	}
}

func (d *Dispatcher) listTools(caller Caller) map[string]any {
	definitions := []map[string]any{}
	for _, tool := range d.tools {
		if !auth.HasScope(caller.Scopes, tool.Scope) {
			continue
		}
		definitions = append(definitions, map[string]any{
			"name":        tool.Name,
			"title":       tool.Title,
			"description": tool.Description,
			"inputSchema": tool.InputSchema,
			"annotations": map[string]any{
				"readOnlyHint":    tool.ReadOnly,
				"destructiveHint": false,
				"idempotentHint":  tool.ReadOnly,
				"openWorldHint":   false,
			},
		})
	}
	return map[string]any{"tools": definitions}
}

func (d *Dispatcher) findTool(name string, caller Caller) (Tool, bool) {
	for _, tool := range d.tools {
		if tool.Name == name && auth.HasScope(caller.Scopes, tool.Scope) {
			return tool, true
		}
	}
	return Tool{}, false
}

func objectSchema(properties map[string]any, required ...string) map[string]any {
	schema := map[string]any{"type": "object", "properties": properties, "additionalProperties": false}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}
