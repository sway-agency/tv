package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
)

func (d *Dispatcher) handleLegacy(ctx context.Context, request Request, envelope Envelope, caller Caller) Outcome {
	if request.Method == "initialize" {
		return d.initialize(request)
	}
	if envelope.HTTP && envelope.ProtocolVersion == ModernVersion {
		return failure(request.ID, http.StatusBadRequest, codeInvalidParams, "missing "+metaVersion+" in _meta", nil)
	}
	if envelope.HTTP && !legacyHeaderAccepted(envelope.ProtocolVersion) {
		data := map[string]any{"supported": supportedVersions(), "requested": envelope.ProtocolVersion}
		return failure(request.ID, http.StatusBadRequest, codeUnsupportedVersion, "Unsupported protocol version", data)
	}
	switch request.Method {
	case "ping":
		return success(request.ID, map[string]any{})
	case "tools/list":
		return success(request.ID, d.listTools(caller))
	case "tools/call":
		return d.callTool(ctx, request, caller)
	}
	return failure(request.ID, http.StatusOK, codeMethodNotFound, "Method not found", nil)
}

func (d *Dispatcher) initialize(request Request) Outcome {
	var params struct {
		ProtocolVersion string `json:"protocolVersion"`
	}
	_ = json.Unmarshal(request.Params, &params)
	version := LegacyVersions[0]
	if slices.Contains(LegacyVersions, params.ProtocolVersion) {
		version = params.ProtocolVersion
	}
	return success(request.ID, map[string]any{
		"protocolVersion": version,
		"capabilities":    map[string]any{"tools": map[string]any{"listChanged": false}},
		"serverInfo":      d.info,
		"instructions":    instructions,
	})
}

func supportedVersions() []string {
	return append([]string{ModernVersion}, LegacyVersions...)
}

func legacyHeaderAccepted(version string) bool {
	return version == "" || slices.Contains(LegacyVersions, version)
}
