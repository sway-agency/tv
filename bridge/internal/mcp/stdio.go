package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const maxMessageBytes = 1 << 20

func ServeStdio(ctx context.Context, input io.Reader, output io.Writer, dispatcher *Dispatcher, caller Caller) error {
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 64<<10), maxMessageBytes)
	encoder := json.NewEncoder(output)
	for scanner.Scan() {
		outcome := handleLine(ctx, scanner.Bytes(), dispatcher, caller)
		if outcome.Response == nil {
			continue
		}
		err := encoder.Encode(outcome.Response)
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}

func handleLine(ctx context.Context, line []byte, dispatcher *Dispatcher, caller Caller) Outcome {
	if len(line) == 0 {
		return Outcome{}
	}
	var request Request
	err := json.Unmarshal(line, &request)
	if err != nil {
		return failure(nil, http.StatusBadRequest, codeParse, "Parse error", nil)
	}
	return dispatcher.Handle(ctx, request, Envelope{}, caller)
}

func ParseError() Outcome {
	return failure(nil, http.StatusBadRequest, codeParse, "Parse error", nil)
}

func InternalError(id json.RawMessage) Outcome {
	return failure(id, http.StatusInternalServerError, codeInternal, "Internal error", nil)
}
