package mcp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"
)

type requestMeta struct {
	Version         string
	HasCapabilities bool
	ToolName        string
}

func readMeta(params json.RawMessage) (requestMeta, error) {
	if len(bytes.TrimSpace(params)) == 0 {
		return requestMeta{}, nil
	}
	var parsed struct {
		Name string                     `json:"name"`
		Meta map[string]json.RawMessage `json:"_meta"`
	}
	err := json.Unmarshal(params, &parsed)
	if err != nil {
		return requestMeta{}, err
	}
	meta := requestMeta{ToolName: parsed.Name}
	rawVersion, found := parsed.Meta[metaVersion]
	if found {
		_ = json.Unmarshal(rawVersion, &meta.Version)
	}
	_, meta.HasCapabilities = parsed.Meta[metaCapabilities]
	return meta, nil
}

func decodeHeaderValue(value string) string {
	encoded, found := strings.CutPrefix(value, "=?base64?")
	if !found || !strings.HasSuffix(encoded, "?=") {
		return value
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSuffix(encoded, "?="))
	if err != nil {
		return "\x00invalid"
	}
	return string(decoded)
}
