package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/aegis-mcp/internal/observability"
)

func TestAuditLogEndToEnd(t *testing.T) {
	// 1. Setup real audit logger with a temp file
	tempDir := t.TempDir()
	auditDBPath := filepath.Join(tempDir, "audit_test.db")
	audit, err := observability.NewAuditLogger(auditDBPath)
	if err != nil {
		t.Fatalf("Failed to create audit logger: %v", err)
	}
	defer audit.Close()

	// 2. Setup backend with the real audit logger
	cfg, engine, readiness, upstream := newTestRuntime(t)
	logger := observability.NewLogger(io.Discard)
	metrics := observability.NewMetrics("test_audit")
	handler := NewServer(cfg, logger, metrics, audit, engine, upstream, readiness...)

	// 3. Make a valid MCP request
	token := signJWT(t, "secret", map[string]any{
		"sub":      "agent-audit-test",
		"agent_id": "agent-audit-test",
		"roles":    []string{"developer"},
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(`{"jsonrpc":"2.0","method":"tools/list","id":1}`))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status OK, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	// 4. Verify audit log entry exists via API
	auditReq := httptest.NewRequest(http.MethodGet, "/api/audit", nil)
	auditRR := httptest.NewRecorder()
	handler.ServeHTTP(auditRR, auditReq)

	if auditRR.Code != http.StatusOK {
		t.Fatalf("Expected status OK for audit log, got %d", auditRR.Code)
	}

	var entries []observability.AuditEntry
	if err := json.Unmarshal(auditRR.Body.Bytes(), &entries); err != nil {
		t.Fatalf("Failed to unmarshal audit entries: %v", err)
	}

	if len(entries) == 0 {
		t.Fatal("Expected at least one audit entry, got zero")
	}

	found := false
	for _, entry := range entries {
		if entry.AgentID == "agent-audit-test" && entry.Method == "tools/list" && entry.Allowed {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Could not find expected audit entry for agent-audit-test and tools/list")
	}
}

func TestAuditLogCapturesDenials(t *testing.T) {
	// 1. Setup real audit logger
	tempDir := t.TempDir()
	auditDBPath := filepath.Join(tempDir, "audit_denial_test.db")
	audit, err := observability.NewAuditLogger(auditDBPath)
	if err != nil {
		t.Fatalf("Failed to create audit logger: %v", err)
	}
	defer audit.Close()

	// 2. Setup backend
	cfg, engine, readiness, upstream := newTestRuntime(t)
	handler := NewServer(cfg, observability.NewLogger(io.Discard), observability.NewMetrics("test_audit"), audit, engine, upstream, readiness...)

	// 3. Make a request that will be blocked by policy (delete_file for developer)
	token := signJWT(t, "secret", map[string]any{
		"sub":      "agent-deny-test",
		"agent_id": "agent-deny-test",
		"roles":    []string{"developer"},
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	// Note: delete_file is only allowed for admin in newTestRuntime's policy
	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(`{"jsonrpc":"2.0","method":"tools/call","params":{"tool":"delete_file"},"id":1}`))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("Expected status Forbidden, got %d", rr.Code)
	}

	// 4. Verify audit log entry reflects the denial
	entries, err := audit.List(10)
	if err != nil {
		t.Fatalf("Failed to list audit entries: %v", err)
	}

	found := false
	for _, entry := range entries {
		if entry.AgentID == "agent-deny-test" && !entry.Allowed && strings.Contains(entry.Reason, "policy") {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Could not find expected denial audit entry")
	}
}
