package runtime

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/aegis-mcp/internal/config"
)

func TestBuildWiresRuntimeDependencies(t *testing.T) {
	t.Parallel()

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
    {"name":"allow-tools-list","effect":"allow","methods":["tools/list"],"roles":["developer"]}
  ]
}`), 0o600); err != nil {
		t.Fatalf("write json file: %v", err)
	}

	deps, err := Build(config.Config{
		ListenAddr:        ":8080",
		UpstreamURL:       "http://127.0.0.1:9090/mcp",
		UpstreamTimeout:   time.Second,
		ReadTimeout:       time.Second,
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      time.Second,
		IdleTimeout:       time.Second,
		MaxBodyBytes:      1024,
		DecisionTimeout:   time.Second,
		MetricsNamespace:  "test",
		Auth: config.AuthConfig{
			SharedSecret: "secret",
		},
		Policy: config.PolicyConfig{
			FilePath:       regoPath,
			ReloadInterval: time.Second,
		},
		Audit: config.AuditConfig{
			FilePath: filepath.Join(t.TempDir(), "audit.db"),
		},
		Semantic: config.SemanticConfig{
			FailClosed:   true,
			MinimumScore: 0.35,
			BlockedIntents: []string{
				"prompt_injection",
			},
		},
	})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if deps.Decision == nil || deps.Policy == nil || deps.Metrics == nil || deps.Logger == nil || deps.Upstream == nil || deps.Inspector == nil {
		t.Fatal("expected all runtime dependencies to be wired")
	}
}
