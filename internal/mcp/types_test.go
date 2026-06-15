package mcp

import (
	"encoding/json"
	"testing"
)

func TestDecodeRequestValid(t *testing.T) {
	t.Parallel()

	raw := []byte(`{"jsonrpc":"2.0","method":"tools/list","id":1}`)
	req, err := DecodeRequest(raw)
	if err != nil {
		t.Fatalf("DecodeRequest() error = %v", err)
	}
	if req.Method != "tools/list" {
		t.Fatalf("method = %q, want %q", req.Method, "tools/list")
	}
}

func TestDecodeRequestRejectsInvalidID(t *testing.T) {
	t.Parallel()

	_, err := DecodeRequest([]byte(`{"jsonrpc":"2.0","method":"tools/list","id":{"nested":true}}`))
	if err == nil {
		t.Fatal("DecodeRequest() expected error")
	}
}

func TestErrorResponsePreservesID(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`"abc"`)
	resp := ErrorResponse(&raw, -32600, "bad request", "detail")
	if resp.ID == nil || string(*resp.ID) != `"abc"` {
		t.Fatalf("ID = %v, want %s", resp.ID, `"abc"`)
	}
}

func FuzzDecodeRequest(f *testing.F) {
	f.Add([]byte(`{"jsonrpc":"2.0","method":"tools/list","id":1}`))
	f.Add([]byte(`{"jsonrpc":"2.0","method":"tools/call","params":{"tool":"read_file"},"id":"req-1"}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`[]`))
	f.Add([]byte(`"string"`))
	f.Add([]byte(` `))

	f.Fuzz(func(t *testing.T, data []byte) {
		req, err := DecodeRequest(data)
		if err == nil {
			// If it parses successfully, it must have a valid method and jsonrpc version
			if req.JSONRPC != JSONRPCVersion {
				t.Errorf("decoded request has invalid jsonrpc version: %q", req.JSONRPC)
			}
			if req.Method == "" {
				t.Errorf("decoded request missing method")
			}
			if req.ID != nil && !validID(*req.ID) {
				t.Errorf("decoded request has invalid ID shape: %s", string(*req.ID))
			}
		}
	})
}
