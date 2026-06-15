package runtime

import (
	"fmt"
	"os"

	"github.com/yourorg/aegis-mcp/internal/auth"
	"github.com/yourorg/aegis-mcp/internal/config"
	"github.com/yourorg/aegis-mcp/internal/decision"
	"github.com/yourorg/aegis-mcp/internal/observability"
	"github.com/yourorg/aegis-mcp/internal/policy"
	"github.com/yourorg/aegis-mcp/internal/proxy"
	"github.com/yourorg/aegis-mcp/internal/semantic"
)

// Dependencies groups the service runtime collaborators.
type Dependencies struct {
	Logger    *observability.Logger
	Metrics   *observability.Metrics
	Audit     *observability.AuditLogger
	Decision  *decision.Engine
	Upstream  proxy.UpstreamClient
	Policy    *policy.Evaluator
	Inspector semantic.Inspector
}

func Build(cfg config.Config) (Dependencies, error) {
	logger := observability.NewLogger(os.Stdout)
	metrics := observability.NewMetrics(cfg.MetricsNamespace)

	auditLogger, err := observability.NewAuditLogger(cfg.Audit.FilePath)
	if err != nil {
		return Dependencies{}, fmt.Errorf("build audit logger: %w", err)
	}

	policyEvaluator, err := policy.NewEvaluator(cfg.Policy)
	if err != nil {
		_ = auditLogger.Close()
		return Dependencies{}, err
	}

	inspector := semantic.NewInspector(cfg.Semantic)
	validator := auth.NewValidator(cfg.Auth)
	engine := decision.NewEngine(validator, inspector, policyEvaluator, cfg.Semantic.MinimumScore, cfg.Semantic.FailClosed, cfg.Semantic.BlockedIntents)
	upstream := proxy.NewHTTPUpstreamClient(cfg.UpstreamURL, cfg.UpstreamTimeout)

	return Dependencies{
		Logger:    logger,
		Metrics:   metrics,
		Audit:     auditLogger,
		Decision:  engine,
		Upstream:  upstream,
		Policy:    policyEvaluator,
		Inspector: inspector,
	}, nil
}
