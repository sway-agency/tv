package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"time"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/mcp"
)

const maxBodyBytes = 64 << 10

type handler struct {
	dispatcher *mcp.Dispatcher
	store      *auth.Store
	limiter    *auth.Limiter
	lockout    *auth.Lockout
	logger     *slog.Logger
}

func (h *handler) serveMCP(writer http.ResponseWriter, request *http.Request) {
	started := time.Now()
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	caller, ok := h.authenticate(writer, request)
	if !ok {
		return
	}
	message, ok := readMessage(writer, request)
	if !ok {
		h.write(writer, mcp.ParseError())
		return
	}
	outcome := h.dispatcher.Handle(request.Context(), message, mcp.Envelope{
		HTTP:            true,
		ProtocolVersion: request.Header.Get("MCP-Protocol-Version"),
		Method:          request.Header.Get("Mcp-Method"),
		Name:            request.Header.Get("Mcp-Name"),
	}, caller)
	h.write(writer, outcome)
	h.audit(request, caller.Name, message.Method, outcome, started)
}

func readMessage(writer http.ResponseWriter, request *http.Request) (mcp.Request, bool) {
	mediaType, _, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		return mcp.Request{}, false
	}
	body, err := io.ReadAll(http.MaxBytesReader(writer, request.Body, maxBodyBytes))
	if err != nil || !bytes.HasPrefix(bytes.TrimSpace(body), []byte("{")) {
		return mcp.Request{}, false
	}
	var message mcp.Request
	err = json.Unmarshal(body, &message)
	return message, err == nil
}

func (h *handler) write(writer http.ResponseWriter, outcome mcp.Outcome) {
	if outcome.Response == nil {
		writer.WriteHeader(outcome.Status)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(outcome.Status)
	_ = json.NewEncoder(writer).Encode(outcome.Response)
}

func (h *handler) audit(request *http.Request, client string, method string, outcome mcp.Outcome, started time.Time) {
	h.logger.Info("request",
		"client", client,
		"address", clientAddress(request),
		"method", method,
		"tool", outcome.Tool,
		"status", outcome.Status,
		"duration_ms", time.Since(started).Milliseconds(),
	)
}
