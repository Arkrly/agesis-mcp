package mcp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

const JSONRPCVersion = "2.0"

// Request models a JSON-RPC request payload used by MCP.
type Request struct {
	JSONRPC string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  *json.RawMessage `json:"params,omitempty"`
	ID      *json.RawMessage `json:"id,omitempty"`
}

// Response models a JSON-RPC response payload.
type Response struct {
	JSONRPC string           `json:"jsonrpc"`
	Result  json.RawMessage  `json:"result,omitempty"`
	Error   *RPCError        `json:"error,omitempty"`
	ID      *json.RawMessage `json:"id,omitempty"`
}

// RPCError is a JSON-RPC error object.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func DecodeRequest(body []byte) (Request, error) {
	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		return Request{}, errors.New("request body is empty")
	}

	var req Request
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		return Request{}, fmt.Errorf("decode request: %w", err)
	}
	if decoder.More() {
		return Request{}, errors.New("request body contains multiple JSON values")
	}
	if err := ValidateRequest(req); err != nil {
		return Request{}, err
	}

	return req, nil
}

func ValidateRequest(req Request) error {
	if req.JSONRPC != JSONRPCVersion {
		return fmt.Errorf("jsonrpc must be %q", JSONRPCVersion)
	}
	if req.Method == "" {
		return errors.New("method is required")
	}
	if req.ID != nil && !validID(*req.ID) {
		return errors.New("id must be a string, number, or null")
	}
	return nil
}

func validID(raw json.RawMessage) bool {
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return false
	}
	switch v.(type) {
	case nil, string, float64:
		return true
	default:
		return false
	}
}

func ErrorResponse(id *json.RawMessage, code int, message string, data any) Response {
	return Response{
		JSONRPC: JSONRPCVersion,
		Error: &RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
		ID: id,
	}
}

// ToolName extracts params.tool from a tools/call request when present.
func (r Request) ToolName() string {
	if r.Params == nil {
		return ""
	}
	var params struct {
		Tool string `json:"tool"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(*r.Params, &params); err != nil {
		return ""
	}
	if params.Tool != "" {
		return params.Tool
	}
	return params.Name
}

// Text returns a conservative string summary of the request for inspection.
func (r Request) Text() string {
	if r.Params == nil {
		return ""
	}
	return string(*r.Params)
}
