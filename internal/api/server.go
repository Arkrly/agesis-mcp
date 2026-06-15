package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"sync"
	"time"

	"github.com/yourorg/aegis-mcp/internal/auth"
	"github.com/yourorg/aegis-mcp/internal/config"
	"github.com/yourorg/aegis-mcp/internal/decision"
	"github.com/yourorg/aegis-mcp/internal/mcp"
	"github.com/yourorg/aegis-mcp/internal/observability"
	"github.com/yourorg/aegis-mcp/internal/policy"
	"github.com/yourorg/aegis-mcp/internal/proxy"
	"golang.org/x/time/rate"
)

const (
	headerSessionID = "Mcp-Session-Id"

	codeParseError     = -32700
	codeInvalidRequest = -32600
	codeUnauthorized   = -32001
	codeForbidden      = -32003
	codeRateLimited    = -32005
	codeInternalError  = -32603
)

// ReadinessChecker reports whether a dependency is ready.
type ReadinessChecker interface {
	Ready() error
}

// Server handles inbound HTTP requests for Aegis-MCP.
type Server struct {
	cfg       config.Config
	logger    *observability.Logger
	metrics   *observability.Metrics
	audit     *observability.AuditLogger
	decision  *decision.Engine
	upstream  proxy.UpstreamClient
	readiness []ReadinessChecker

	mu           sync.Mutex
	limiters     map[string]*rate.Limiter
	limitContext context.Context
}

func NewServer(
	cfg config.Config,
	logger *observability.Logger,
	metrics *observability.Metrics,
	audit *observability.AuditLogger,
	engine *decision.Engine,
	upstream proxy.UpstreamClient,
	readiness ...ReadinessChecker,
) http.Handler {
	server := &Server{
		cfg:          cfg,
		logger:       logger,
		metrics:      metrics,
		audit:        audit,
		decision:     engine,
		upstream:     upstream,
		readiness:    readiness,
		limiters:     make(map[string]*rate.Limiter),
		limitContext: context.Background(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /live", server.handleLive)
	mux.HandleFunc("GET /ready", server.handleReady)
	mux.Handle("GET /metrics", server.metrics.Handler())
	mux.HandleFunc("POST /mcp", server.handleMCP)
	mux.HandleFunc("GET /api/audit", server.handleAudit)
	mux.HandleFunc("GET /api/summary", server.handleSummary)

	return server.withMiddleware(mux)
}

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	// In a real app, we'd check for 'auditor' role here
	limit := 100
	entries, err := s.audit.List(limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

func (s *Server) handleSummary(w http.ResponseWriter, _ *http.Request) {
	// Basic health summary for the dashboard
	status := "ok"
	for _, checker := range s.readiness {
		if err := checker.Ready(); err != nil {
			status = "degraded"
			break
		}
	}

	summary := map[string]any{
		"status":    status,
		"version":   "v0.1.0",
		"timestamp": time.Now().UTC(),
		"components": map[string]string{
			"proxy":     "healthy",
			"policy":    "healthy",
			"audit_log": "healthy",
		},
	}
	writeJSON(w, http.StatusOK, summary)
}

func (s *Server) withMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		status := http.StatusOK

		defer func() {
			if recovered := recover(); recovered != nil {
				status = http.StatusInternalServerError
				s.logger.Error("panic recovered", map[string]any{
					"path":  r.URL.Path,
					"error": recovered,
				})
				writeJSON(w, http.StatusInternalServerError, mcp.ErrorResponse(nil, codeInternalError, "internal server error", nil))
			}
			s.metrics.ObserveRequest(r.URL.Path, status, time.Since(start))
			s.logger.Info("request completed", map[string]any{
				"method":   r.Method,
				"path":     r.URL.Path,
				"remote":   r.RemoteAddr,
				"status":   status,
				"duration": time.Since(start).String(),
			})
		}()

		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		status = recorder.status
	})
}

func (s *Server) handleLive(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleReady(w http.ResponseWriter, _ *http.Request) {
	for _, checker := range s.readiness {
		if err := checker.Ready(); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{
				"status": "not_ready",
				"error":  err.Error(),
			})
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (s *Server) handleMCP(w http.ResponseWriter, r *http.Request) {
	if !contentTypeIsJSON(r.Header.Get("Content-Type")) {
		writeJSON(w, http.StatusUnsupportedMediaType, mcp.ErrorResponse(nil, codeInvalidRequest, "content type must be application/json", nil))
		return
	}

	body, err := readRequestBody(http.MaxBytesReader(w, r.Body, s.cfg.MaxBodyBytes))
	if err != nil {
		status := http.StatusBadRequest
		message := "invalid request body"
		if errors.As(err, new(*http.MaxBytesError)) {
			status = http.StatusRequestEntityTooLarge
			message = "request body exceeds configured limit"
		}
		s.logger.Error("request body rejected", map[string]any{"error": err.Error()})
		writeJSON(w, status, mcp.ErrorResponse(nil, codeParseError, message, nil))
		return
	}

	req, err := mcp.DecodeRequest(body)
	if err != nil {
		s.logger.Error("mcp decode rejected", map[string]any{"error": err.Error()})
		writeJSON(w, http.StatusBadRequest, mcp.ErrorResponse(nil, codeInvalidRequest, "invalid MCP request", err.Error()))
		return
	}

	token, err := auth.ExtractBearerToken(r.Header.Get("Authorization"))
	if err != nil {
		s.metrics.ObserveDenial("auth")
		writeJSON(w, http.StatusUnauthorized, mcp.ErrorResponse(req.ID, codeUnauthorized, "missing or invalid authorization", nil))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), s.cfg.DecisionTimeout)
	defer cancel()

	result, err := s.decision.Evaluate(ctx, token, req)

	// Per-agent rate limiting
	if err == nil {
		if !s.allowRate(result.Claims.AgentID) {
			s.metrics.ObserveDenial("rate_limit")
			writeJSON(w, http.StatusTooManyRequests, mcp.ErrorResponse(req.ID, codeRateLimited, "rate limit exceeded", nil))
			return
		}
	}

	// Audit Log
	if s.audit != nil {
		_ = s.audit.Log(observability.AuditEntry{
			AgentID: result.Claims.AgentID,
			Method:  req.Method,
			Tool:    req.ToolName(),
			Allowed: err == nil,
			Reason:  func() string { if err != nil { return err.Error() }; return "" }(),
			Inspection: map[string]any{
				"safety_score":      result.Inspection.SafetyScore,
				"intent_categories": result.Inspection.IntentCategories,
			},
		})
	}

	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidAuthorization),
			errors.Is(err, auth.ErrInvalidToken),
			errors.Is(err, auth.ErrExpiredToken),
			errors.Is(err, auth.ErrTokenNotYetValid),
			errors.Is(err, auth.ErrMissingAgentID):
			s.metrics.ObserveDenial("auth")
			writeJSON(w, http.StatusUnauthorized, mcp.ErrorResponse(req.ID, codeUnauthorized, "authentication failed", map[string]any{"reason": err.Error()}))
			return
		case errors.Is(err, decision.ErrSemanticDenied):
			s.metrics.ObserveDenial("semantic")
			writeJSON(w, http.StatusUnprocessableEntity, mcp.ErrorResponse(req.ID, codeForbidden, "request blocked by semantic inspection", map[string]any{
				"agent_id":          result.Claims.AgentID,
				"safety_score":      result.Inspection.SafetyScore,
				"intent_categories": result.Inspection.IntentCategories,
			}))
			return
		case errors.Is(err, policy.ErrPolicyDenied):
			s.metrics.ObserveDenial("policy")
			writeJSON(w, http.StatusForbidden, mcp.ErrorResponse(req.ID, codeForbidden, "request blocked by policy", map[string]any{
				"agent_id": result.Claims.AgentID,
				"policy":   result.Policy.PolicyName,
				"reason":   result.Policy.Reason,
				"metadata": result.Policy.Metadata,
			}))
			return
		default:
			s.metrics.ObserveDenial("internal")
			s.logger.Error("decision engine failed", map[string]any{"error": err.Error()})
			writeJSON(w, http.StatusInternalServerError, mcp.ErrorResponse(req.ID, codeInternalError, "decision engine failure", nil))
			return
		}
	}

	resp, err := s.upstream.Forward(r.Context(), body, r.Header.Get("Authorization"), r.Header.Get(headerSessionID))
	if err != nil {
		s.logger.Error("upstream request failed", map[string]any{
			"method": req.Method,
			"tool":   req.ToolName(),
			"error":  err.Error(),
		})
		writeJSON(w, http.StatusBadGateway, mcp.ErrorResponse(req.ID, codeInternalError, "upstream request failed", nil))
		return
	}

	if err := proxy.CopyResponse(w, resp); err != nil {
		s.logger.Error("copy upstream response failed", map[string]any{"error": err.Error()})
	}
}

func (s *Server) allowRate(agentID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	limiter, ok := s.limiters[agentID]
	if !ok {
		// Default: 10 requests per second, burst of 20
		limiter = rate.NewLimiter(rate.Limit(10), 20)
		s.limiters[agentID] = limiter
	}

	return limiter.Allow()
}

func readRequestBody(body io.ReadCloser) ([]byte, error) {
	defer body.Close()

	payload, err := io.ReadAll(body)
	if errors.Is(err, io.EOF) {
		return payload, nil
	}
	return payload, err
}

func contentTypeIsJSON(value string) bool {
	if value == "" {
		return false
	}
	mediaType, _, err := mime.ParseMediaType(value)
	if err != nil {
		return false
	}
	return mediaType == "application/json"
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, `{"jsonrpc":"2.0","error":{"code":-32603,"message":"internal server error"}}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
