package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      any             `json:"id,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  any    `json:"result,omitempty"`
	Error   any    `json:"error,omitempty"`
	ID      any    `json:"id,omitempty"`
}

func main() {
	port := 9090
	mux := http.NewServeMux()

	mux.HandleFunc("/mcp", handleMCP)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Mock MCP Server is running on :%d/mcp", port)
	})

	log.Printf("Starting Mock MCP Server on :%d...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatal(err)
	}
}

func handleMCP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	var req JSONRPCRequest
	if err := json.Unmarshal(body, &req); err != nil {
		sendError(w, nil, -32700, "Parse error")
		return
	}

	log.Printf("Received MCP Request: %s", req.Method)

	switch req.Method {
	case "tools/list":
		handleToolsList(w, req.ID)
	case "tools/call":
		handleToolsCall(w, req.ID, req.Params)
	default:
		sendError(w, req.ID, -32601, "Method not found")
	}
}

func handleToolsList(w http.ResponseWriter, id any) {
	tools := []map[string]any{
		{
			"name":        "echo",
			"description": "Repeat back the input string",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"message": map[string]any{"type": "string"},
				},
				"required": []string{"message"},
			},
		},
		{
			"name":        "get_system_info",
			"description": "Get mock system status and metrics",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{},
			},
		},
	}

	sendResponse(w, id, map[string]any{
		"tools": tools,
	})
}

func handleToolsCall(w http.ResponseWriter, id any, params json.RawMessage) {
	var p struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		sendError(w, id, -32602, "Invalid params")
		return
	}

	switch p.Name {
	case "echo":
		msg := p.Arguments["message"]
		sendResponse(w, id, map[string]any{
			"content": []map[string]any{
				{"type": "text", "text": fmt.Sprintf("Echo: %v", msg)},
			},
		})
	case "get_system_info":
		sendResponse(w, id, map[string]any{
			"content": []map[string]any{
				{"type": "text", "text": "System Status: Healthy\nUptime: 14h 22m\nLoad: 0.12, 0.08, 0.05"},
			},
		})
	default:
		sendError(w, id, -32602, "Tool not found")
	}
}

func sendResponse(w http.ResponseWriter, id any, result any) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func sendError(w http.ResponseWriter, id any, code int, message string) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		Error: map[string]any{
			"code":    code,
			"message": message,
		},
		ID: id,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // JSON-RPC usually returns 200 even for RPC errors
	json.NewEncoder(w).Encode(resp)
}
