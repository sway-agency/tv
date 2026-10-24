package server

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"time"

	"github.com/sway-agency/tv/bridge/internal/mcp"
)

func (h *handler) serveTool(writer http.ResponseWriter, request *http.Request) {
	started := time.Now()
	caller, ok := h.authenticate(writer, request)
	if !ok {
		return
	}
	arguments, ok := readArguments(writer, request)
	if !ok {
		http.Error(writer, "body must be a JSON object of tool arguments", http.StatusBadRequest)
		return
	}
	params, err := json.Marshal(map[string]any{"name": request.PathValue("name"), "arguments": arguments})
	if err != nil {
		h.write(writer, mcp.InternalError(nil))
		return
	}
	message := mcp.Request{JSONRPC: "2.0", ID: json.RawMessage("1"), Method: "tools/call", Params: params}
	outcome := h.dispatcher.Handle(request.Context(), message, mcp.Envelope{}, caller)
	if outcome.Response != nil && outcome.Response.Error == nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(outcome.Status)
		_ = json.NewEncoder(writer).Encode(outcome.Response.Result)
	} else {
		h.write(writer, outcome)
	}
	h.audit(request, caller.Name, "rest", outcome, started)
}

func readArguments(writer http.ResponseWriter, request *http.Request) (json.RawMessage, bool) {
	mediaType, _, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		return nil, false
	}
	body, err := io.ReadAll(http.MaxBytesReader(writer, request.Body, maxBodyBytes))
	if err != nil {
		return nil, false
	}
	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		return json.RawMessage("{}"), true
	}
	return body, bytes.HasPrefix(body, []byte("{")) && json.Valid(body)
}
