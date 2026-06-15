package policy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage/inmem"
	"github.com/yourorg/aegis-mcp/internal/auth"
	"github.com/yourorg/aegis-mcp/internal/config"
	"github.com/yourorg/aegis-mcp/internal/mcp"
	"github.com/yourorg/aegis-mcp/internal/semantic"
)

var ErrPolicyDenied = errors.New("policy denied request")

// Decision is the policy evaluation result.
type Decision struct {
	Allowed    bool           `json:"allowed"`
	PolicyName string         `json:"policy_name,omitempty"`
	Reason     string         `json:"reason,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

// Evaluator checks requests against the loaded policy document using OPA.
type Evaluator struct {
	path          string
	reloadEvery   time.Duration
	now           func() time.Time
	mu            sync.RWMutex
	query         rego.PreparedEvalQuery
	cache         *lru.Cache[string, Decision]
	lastLoaded    time.Time
	lastModTime   time.Time
	lastLoadError error
}

func NewEvaluator(cfg config.PolicyConfig) (*Evaluator, error) {
	cache, err := lru.New[string, Decision](1000)
	if err != nil {
		return nil, fmt.Errorf("create decision cache: %w", err)
	}

	e := &Evaluator{
		path:        cfg.FilePath,
		reloadEvery: cfg.ReloadInterval,
		now:         time.Now,
		cache:       cache,
	}
	if err := e.Reload(); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Evaluator) Ready() error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.lastLoadError
}

func (e *Evaluator) Reload() error {
	regoCode, err := os.ReadFile(e.path)
	if err != nil {
		return fmt.Errorf("read rego file: %w", err)
	}

	// Load data file (e.g., policies.json) if it exists
	var data map[string]any
	dataPath := strings.TrimSuffix(e.path, filepath.Ext(e.path)) + ".json"
	if _, err := os.Stat(dataPath); err == nil {
		dataBytes, err := os.ReadFile(dataPath)
		if err != nil {
			return fmt.Errorf("read policy data file: %w", err)
		}
		if err := json.Unmarshal(dataBytes, &data); err != nil {
			return fmt.Errorf("unmarshal policy data: %w", err)
		}
	}

	ctx := context.Background()
	r := rego.New(
		rego.Query("data.aegis.mcp.decision"),
		rego.Module("policy.rego", string(regoCode)),
	)

	if data != nil {
		r = rego.New(
			rego.Query("data.aegis.mcp.decision"),
			rego.Module("policy.rego", string(regoCode)),
			rego.Store(inmem.NewFromObject(data)),
		)
	}

	query, err := r.PrepareForEval(ctx)
	if err != nil {
		return fmt.Errorf("prepare rego query: %w", err)
	}

	info, err := os.Stat(e.path)
	if err != nil {
		return fmt.Errorf("stat rego file: %w", err)
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	e.query = query
	e.lastLoaded = e.now()
	e.lastModTime = info.ModTime()
	e.lastLoadError = nil
	e.cache.Purge() // Clear cache on reload
	return nil
}

func (e *Evaluator) Evaluate(ctx context.Context, req mcp.Request, claims auth.Claims, inspection semantic.Result) (Decision, error) {
	if err := e.reloadIfNeeded(); err != nil {
		return Decision{}, err
	}

	// Check cache
	cacheKey := fmt.Sprintf("%s:%s:%s:%f", claims.AgentID, req.Method, req.ToolName(), inspection.SafetyScore)
	if decision, ok := e.cache.Get(cacheKey); ok {
		if !decision.Allowed {
			return decision, ErrPolicyDenied
		}
		return decision, nil
	}

	e.mu.RLock()
	query := e.query
	e.mu.RUnlock()

	input := map[string]any{
		"method": req.Method,
		"tool":   req.ToolName(),
		"auth": map[string]any{
			"agent_id": claims.AgentID,
			"roles":    claims.Roles,
		},
		"inspection": map[string]any{
			"safety_score":      inspection.SafetyScore,
			"intent_categories": inspection.IntentCategories,
		},
	}

	results, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return Decision{}, fmt.Errorf("eval policy: %w", err)
	}

	if len(results) == 0 {
		return Decision{Allowed: false, Reason: "no policy results"}, ErrPolicyDenied
	}

	decisionMap, ok := results[0].Expressions[0].Value.(map[string]any)
	if !ok {
		return Decision{Allowed: false, Reason: "invalid policy output shape"}, ErrPolicyDenied
	}

	allowed, _ := decisionMap["final_allow"].(bool)
	decision := Decision{
		Allowed: allowed,
	}

	if !allowed {
		decision.Reason = "blocked by OPA policy"
		if denied, _ := decisionMap["denied"].(bool); denied {
			decision.Reason = "explicitly denied by OPA policy"
		}
		e.cache.Add(cacheKey, decision)
		return decision, ErrPolicyDenied
	}

	e.cache.Add(cacheKey, decision)
	return decision, nil
}

func (e *Evaluator) reloadIfNeeded() error {
	e.mu.RLock()
	path := e.path
	lastLoaded := e.lastLoaded
	interval := e.reloadEvery
	e.mu.RUnlock()

	if interval > 0 && e.now().Sub(lastLoaded) < interval {
		return nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	e.mu.RLock()
	lastModTime := e.lastModTime
	e.mu.RUnlock()
	if !info.ModTime().After(lastModTime) {
		e.mu.Lock()
		e.lastLoaded = e.now()
		e.mu.Unlock()
		return nil
	}

	return e.Reload()
}
