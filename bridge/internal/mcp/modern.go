package mcp

import (
	"context"
	"net/http"
)

func (d *Dispatcher) handleModern(ctx context.Context, request Request, envelope Envelope, caller Caller, meta requestMeta) Outcome {
	if envelope.HTTP {
		mismatch := headerMismatch(request, envelope, meta)
		if mismatch != "" {
			return failure(request.ID, http.StatusBadRequest, codeHeaderMismatch, mismatch, nil)
		}
	}
	if meta.Version != ModernVersion {
		data := map[string]any{"supported": supportedVersions(), "requested": meta.Version}
		return failure(request.ID, http.StatusBadRequest, codeUnsupportedVersion, "Unsupported protocol version", data)
	}
	if !meta.HasCapabilities {
		return failure(request.ID, http.StatusBadRequest, codeInvalidParams, "missing "+metaCapabilities+" in _meta", nil)
	}
	switch request.Method {
	case "server/discover":
		return success(request.ID, d.modernResult(map[string]any{
			"supportedVersions": supportedVersions(),
			"capabilities":      map[string]any{"tools": map[string]any{}},
			"instructions":      instructions,
		}))
	case "ping":
		return success(request.ID, d.modernResult(map[string]any{}))
	case "tools/list":
		return success(request.ID, d.modernResult(d.listTools(caller)))
	case "tools/call":
		return d.modernCall(ctx, request, caller)
	}
	return failure(request.ID, http.StatusNotFound, codeMethodNotFound, "Method not found", nil)
}

func (d *Dispatcher) modernCall(ctx context.Context, request Request, caller Caller) Outcome {
	outcome := d.callTool(ctx, request, caller)
	if outcome.Response != nil && outcome.Response.Error == nil {
		result, _ := outcome.Response.Result.(map[string]any)
		outcome.Response.Result = d.modernResult(result)
	}
	return outcome
}

func (d *Dispatcher) modernResult(result map[string]any) map[string]any {
	result["resultType"] = "complete"
	result["_meta"] = map[string]any{metaServerInfo: d.info}
	return result
}

func headerMismatch(request Request, envelope Envelope, meta requestMeta) string {
	if envelope.ProtocolVersion != meta.Version {
		return "Header mismatch: MCP-Protocol-Version does not match _meta protocolVersion"
	}
	if envelope.Method != request.Method {
		return "Header mismatch: Mcp-Method does not match the request method"
	}
	if request.Method == "tools/call" && decodeHeaderValue(envelope.Name) != meta.ToolName {
		return "Header mismatch: Mcp-Name does not match params.name"
	}
	return ""
}
