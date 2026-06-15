package policy

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/aegis-mcp/internal/auth"
	"github.com/yourorg/aegis-mcp/internal/config"
	"github.com/yourorg/aegis-mcp/internal/mcp"
	"github.com/yourorg/aegis-mcp/internal/semantic"
)

func TestEvaluatorAllowsMatchingRule(t *testing.T) {
	t.Parallel()

	evaluator := newEvaluator(t)
	req := mcp.Request{JSONRPC: "2.0", Method: "tools/list"}
	claims := auth.Claims{AgentID: "agent-1", Roles: []string{"developer"}}
	inspection := semantic.Result{SafetyScore: 1}

	decision, err := evaluator.Evaluate(context.Background(), req, claims, inspection)
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if !decision.Allowed {
		t.Fatal("expected allow decision")
	}
}

func TestEvaluatorRejectsBlockedIntent(t *testing.T) {
	t.Parallel()

	evaluator := newEvaluator(t)
	req := mcp.Request{JSONRPC: "2.0", Method: "tools/list"}
	claims := auth.Claims{AgentID: "agent-1", Roles: []string{"developer"}}
	inspection := semantic.Result{SafetyScore: 1, IntentCategories: []string{"prompt_injection"}}

	_, err := evaluator.Evaluate(context.Background(), req, claims, inspection)
	if err != ErrPolicyDenied {
		t.Fatalf("Evaluate() error = %v, want %v", err, ErrPolicyDenied)
	}
}

func TestEvaluatorRejectsMissingRule(t *testing.T) {
	t.Parallel()

	evaluator := newEvaluator(t)
	raw := json.RawMessage(`{"tool":"delete_file"}`)
	req := mcp.Request{JSONRPC: "2.0", Method: "tools/call", Params: &raw}
	claims := auth.Claims{AgentID: "agent-1", Roles: []string{"developer"}}
	inspection := semantic.Result{SafetyScore: 1}

	_, err := evaluator.Evaluate(context.Background(), req, claims, inspection)
	if err != ErrPolicyDenied {
		t.Fatalf("Evaluate() error = %v, want %v", err, ErrPolicyDenied)
	}
}

func newEvaluator(t *testing.T) *Evaluator {
	t.Helper()

	basePath := filepath.Join(t.TempDir(), "policy")
	regoPath := basePath + ".rego"
	jsonPath := basePath + ".json"

	if err := os.WriteFile(regoPath, []byte(`package aegis.mcp
import rego.v1
default allow = false
allow if {
    some rule in data.rules
    rule.effect == "allow"
    rule.methods[_] == input.method
    input.inspection.safety_score >= data.semantic.minimum_safety_score
    check_intents(input.inspection.intent_categories)
}
check_intents(intents) if {
    count({i | i := intents[_]; contains_blocked(i)}) == 0
}
contains_blocked(i) if {
    some blocked in data.semantic.blocked_intents
    blocked == i
}
decision := {"final_allow": allow}
`), 0o600); err != nil {
		t.Fatalf("write rego file: %v", err)
	}

	if err := os.WriteFile(jsonPath, []byte(`{
  "semantic": {
    "minimum_safety_score": 0.35,
    "blocked_intents": ["prompt_injection"]
  },
  "rules": [
    {"name":"allow-tools-list","effect":"allow","methods":["tools/list"]}
  ]
}`), 0o600); err != nil {
		t.Fatalf("write json file: %v", err)
	}

	evaluator, err := NewEvaluator(config.PolicyConfig{
		FilePath:       regoPath,
		ReloadInterval: time.Second,
	})
	if err != nil {
		t.Fatalf("NewEvaluator() error = %v", err)
	}
	return evaluator
}
