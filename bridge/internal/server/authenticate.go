package server

import (
	"net/http"

	"github.com/sway-agency/tv/bridge/internal/mcp"
)

func (h *handler) authenticate(writer http.ResponseWriter, request *http.Request) (mcp.Caller, bool) {
	address := lockoutKey(clientAddress(request))
	if h.lockout.Locked(address) {
		h.logger.Warn("locked out", "address", address)
		http.Error(writer, "too many failed attempts", http.StatusTooManyRequests)
		return mcp.Caller{}, false
	}
	token, ok := h.store.Authenticate(request.Header.Get("Authorization"))
	if !ok {
		h.lockout.Fail(address)
		h.logger.Warn("authentication failed", "address", address)
		writer.Header().Set("WWW-Authenticate", `Bearer realm="accueil-bridge"`)
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		return mcp.Caller{}, false
	}
	if !h.limiter.Allow(token.Name) {
		writer.Header().Set("Retry-After", "1")
		http.Error(writer, "rate limit exceeded", http.StatusTooManyRequests)
		return mcp.Caller{}, false
	}
	return mcp.Caller{Name: token.Name, Scopes: token.Scopes}, true
}
