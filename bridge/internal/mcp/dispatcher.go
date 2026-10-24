package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/tv"
)

const instructions = "Controls a living-room Android TV running the Accueil launcher. " +
	"Use tv_status first, prefer tv_launch_app over key navigation, and ask the user before turning the TV off."

type Caller struct {
	Name   string
	Scopes []auth.Scope
}

type Dispatcher struct {
	controller tv.Controller
	tools      []Tool
	info       Implementation
}

func NewDispatcher(controller tv.Controller, version string) *Dispatcher {
	return &Dispatcher{controller: controller, tools: DefaultTools(), info: Implementation{Name: "accueil-bridge", Version: version}}
}

func (d *Dispatcher) Handle(ctx context.Context, request Request, envelope Envelope, caller Caller) Outcome {
	if request.JSONRPC != "2.0" || request.Method == "" {
		return failure(nil, http.StatusBadRequest, codeInvalidRequest, "invalid JSON-RPC request", nil)
	}
	if len(request.ID) == 0 {
		return Outcome{Status: http.StatusAccepted}
	}
	if !validID(request.ID) {
		return failure(nil, http.StatusBadRequest, codeInvalidRequest, "request id must be a string or a number", nil)
	}
	meta, err := readMeta(request.Params)
	if err != nil {
		return failure(request.ID, http.StatusBadRequest, codeInvalidParams, "params must be an object", nil)
	}
	if meta.Version != "" {
		return d.handleModern(ctx, request, envelope, caller, meta)
	}
	return d.handleLegacy(ctx, request, envelope, caller)
}

func validID(id json.RawMessage) bool {
	trimmed := bytes.TrimSpace(id)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return false
	}
	var value any
	err := json.Unmarshal(trimmed, &value)
	if err != nil {
		return false
	}
	switch value.(type) {
	case string, float64:
		return true
	}
	return false
}
