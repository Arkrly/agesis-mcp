package semantic

import (
	"context"
	"strings"

	"github.com/yourorg/aegis-mcp/internal/config"
	"github.com/yourorg/aegis-mcp/internal/mcp"
)

// Result captures semantic inspection output.
type Result struct {
	SafetyScore      float64        `json:"safety_score"`
	IntentCategories []string       `json:"intent_categories"`
	Explanation      string         `json:"explanation,omitempty"`
	Metadata         map[string]any `json:"metadata,omitempty"`
}

// Inspector analyzes requests for safety and intent using a specific provider.
type Inspector interface {
	Inspect(context.Context, mcp.Request) (Result, error)
	Ready() error
}

// NewInspector returns an inspector based on the configuration.
func NewInspector(cfg config.SemanticConfig) Inspector {
	// For now, only heuristic is supported in this build.
	// Future: support llama.cpp or external safety APIs.
	return NewHeuristicInspector(cfg)
}

type HeuristicInspector struct {
	cfg config.SemanticConfig
}

func NewHeuristicInspector(cfg config.SemanticConfig) *HeuristicInspector {
	return &HeuristicInspector{cfg: cfg}
}

func (i *HeuristicInspector) Ready() error {
	return nil
}

func (i *HeuristicInspector) Inspect(_ context.Context, req mcp.Request) (Result, error) {
	summary := strings.ToLower(req.Method + " " + req.Text())
	intents := make([]string, 0, 4)
	score := 1.0

	if containsAny(summary, "ignore previous instructions", "ignore all previous", "system prompt", "reveal your instructions") {
		intents = append(intents, "prompt_injection")
		score -= 0.55
	}
	if containsAny(summary, "jailbreak", "developer mode", "bypass safeguards", "disable safety") {
		intents = append(intents, "jailbreak_attempt")
		score -= 0.45
	}
	if containsAny(summary, "ransomware", "keylogger", "credential theft", "steal secrets", "exfiltrate", "delete all files") {
		intents = append(intents, "malicious_code_generation")
		score -= 0.70
	}
	if containsAny(summary, "password", "api key", "token", "secret", "credential") && containsAny(summary, "dump", "print", "expose", "send", "upload") {
		intents = append(intents, "secret_exfiltration")
		score -= 0.60
	}

	if req.Method == "tools/call" {
		tool := strings.ToLower(req.ToolName())
		switch tool {
		case "delete_file", "remove_file", "exec_command", "run_shell":
			score -= 0.10
		}
	}

	if score < 0 {
		score = 0
	}

	return Result{
		SafetyScore:      score,
		IntentCategories: dedupe(intents),
		Metadata: map[string]any{
			"mode": "heuristic",
		},
	}, nil
}

func containsAny(haystack string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(haystack, needle) {
			return true
		}
	}
	return false
}

func dedupe(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
