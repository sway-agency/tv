package mcp

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/sway-agency/tv/bridge/internal/tv"
)

const toolTimeout = 30 * time.Second

func (d *Dispatcher) callTool(ctx context.Context, request Request, caller Caller) Outcome {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	err := json.Unmarshal(request.Params, &params)
	if err != nil {
		return failure(request.ID, http.StatusBadRequest, codeInvalidParams, "invalid tools/call params", nil)
	}
	tool, found := d.findTool(params.Name, caller)
	if !found {
		return failure(request.ID, http.StatusOK, codeInvalidParams, "Unknown tool: "+params.Name, nil)
	}
	ctx, cancel := context.WithTimeout(ctx, toolTimeout)
	defer cancel()
	result, err := tool.Run(ctx, d.controller, params.Arguments)
	outcome := success(request.ID, toolResult(result, err))
	outcome.Tool = tool.Name
	return outcome
}

func toolResult(result ToolResult, err error) map[string]any {
	if err != nil {
		message := err.Error()
		if errors.Is(err, tv.ErrCommandFailed) {
			message = "The TV did not accept the command. Check that it is on and reachable."
		}
		return map[string]any{"content": []any{textContent(message)}, "isError": true}
	}
	content := []any{}
	if result.Text != "" {
		content = append(content, textContent(result.Text))
	}
	if result.ImageJPEG != nil {
		content = append(content, map[string]any{
			"type": "image", "mimeType": "image/jpeg", "data": base64.StdEncoding.EncodeToString(result.ImageJPEG),
		})
	}
	response := map[string]any{"content": content, "isError": false}
	if result.Structured != nil {
		response["structuredContent"] = result.Structured
	}
	return response
}

func textContent(text string) map[string]any {
	return map[string]any{"type": "text", "text": text}
}

func decodeArguments(raw json.RawMessage, target any) error {
	if len(bytes.TrimSpace(raw)) == 0 || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		raw = json.RawMessage("{}")
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	err := decoder.Decode(target)
	if err != nil {
		return errors.New("invalid arguments: " + err.Error())
	}
	return nil
}
