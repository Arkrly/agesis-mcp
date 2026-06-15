package decision

import (
	"context"
	"errors"
	"fmt"

	"github.com/yourorg/aegis-mcp/internal/auth"
	"github.com/yourorg/aegis-mcp/internal/mcp"
	"github.com/yourorg/aegis-mcp/internal/policy"
	"github.com/yourorg/aegis-mcp/internal/semantic"
)

var ErrSemanticDenied = errors.New("semantic inspection denied request")

// Engine combines auth, semantic inspection, and policy evaluation.
type Engine struct {
	auth           *auth.Validator
	inspector      semantic.Inspector
	policy         *policy.Evaluator
	minScore       float64
	failClosed     bool
	blockedIntents map[string]struct{}
}

type Result struct {
	Claims     auth.Claims     `json:"claims"`
	Inspection semantic.Result `json:"inspection"`
	Policy     policy.Decision `json:"policy"`
}

func NewEngine(authValidator *auth.Validator, inspector semantic.Inspector, evaluator *policy.Evaluator, minScore float64, failClosed bool, blockedIntents []string) *Engine {
	intentSet := make(map[string]struct{}, len(blockedIntents))
	for _, intent := range blockedIntents {
		intentSet[intent] = struct{}{}
	}
	return &Engine{
		auth:           authValidator,
		inspector:      inspector,
		policy:         evaluator,
		minScore:       minScore,
		failClosed:     failClosed,
		blockedIntents: intentSet,
	}
}

func (e *Engine) Evaluate(ctx context.Context, bearerToken string, req mcp.Request) (Result, error) {
	claims, err := e.auth.Validate(bearerToken)
	if err != nil {
		return Result{}, err
	}

	inspection, err := e.inspector.Inspect(ctx, req)
	if err != nil {
		if e.failClosed {
			return Result{}, fmt.Errorf("semantic inspection failed: %w", err)
		}
		inspection = semantic.Result{
			SafetyScore: 1,
			Metadata:    map[string]any{"degraded": true},
		}
	}

	if inspection.SafetyScore < e.minScore {
		return Result{
			Claims:     claims,
			Inspection: inspection,
		}, ErrSemanticDenied
	}
	for _, intent := range inspection.IntentCategories {
		if _, blocked := e.blockedIntents[intent]; blocked {
			return Result{
				Claims:     claims,
				Inspection: inspection,
			}, ErrSemanticDenied
		}
	}

	decision, err := e.policy.Evaluate(ctx, req, claims, inspection)
	result := Result{
		Claims:     claims,
		Inspection: inspection,
		Policy:     decision,
	}
	if err != nil {
		return result, err
	}
	return result, nil
}
