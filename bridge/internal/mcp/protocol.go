package mcp

import (
	"encoding/json"
	"net/http"
)

const (
	ModernVersion    = "2026-07-28"
	legacyDefault    = "2025-03-26"
	metaVersion      = "io.modelcontextprotocol/protocolVersion"
	metaCapabilities = "io.modelcontextprotocol/clientCapabilities"
	metaServerInfo   = "io.modelcontextprotocol/serverInfo"
)

var LegacyVersions = []string{"2025-11-25", "2025-06-18", "2025-03-26"}

const (
	codeParse              = -32700
	codeInvalidRequest     = -32600
	codeMethodNotFound     = -32601
	codeInvalidParams      = -32602
	codeInternal           = -32603
	codeHeaderMismatch     = -32020
	codeUnsupportedVersion = -32022
)

type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type Implementation struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Envelope struct {
	HTTP            bool
	ProtocolVersion string
	Method          string
	Name            string
}

type Outcome struct {
	Response *Response
	Status   int
	Tool     string
}

func failure(id json.RawMessage, status int, code int, message string, data any) Outcome {
	return Outcome{
		Response: &Response{JSONRPC: "2.0", ID: id, Error: &Error{Code: code, Message: message, Data: data}},
		Status:   status,
	}
}

func success(id json.RawMessage, result any) Outcome {
	return Outcome{Response: &Response{JSONRPC: "2.0", ID: id, Result: result}, Status: http.StatusOK}
}
