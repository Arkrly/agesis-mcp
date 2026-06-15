package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultListenAddr           = ":8080"
	defaultReadTimeout          = 10 * time.Second
	defaultReadHeaderTimeout    = 5 * time.Second
	defaultWriteTimeout         = 15 * time.Second
	defaultIdleTimeout          = 60 * time.Second
	defaultUpstreamTimeout      = 30 * time.Second
	defaultMaxBodyBytes         = 1 << 20
	defaultDecisionTimeout      = 5 * time.Second
	defaultSemanticMinimumScore = 0.35
	defaultMetricsNamespace     = "aegis_mcp"
	defaultSemanticFailClosed   = true
	defaultPolicyReloadInterval = 5 * time.Second
)

// Config stores process-level HTTP proxy configuration.
type Config struct {
	ListenAddr        string
	UpstreamURL       string
	UpstreamTimeout   time.Duration
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxBodyBytes      int64
	DecisionTimeout   time.Duration
	MetricsNamespace  string

	Auth     AuthConfig
	Policy   PolicyConfig
	Semantic SemanticConfig
	Audit    AuditConfig
}

// AuthConfig stores JWT validation settings.
type AuthConfig struct {
	SharedSecret     string
	Issuer           string
	Audience         []string
	RequireAudience  bool
	AllowedClockSkew time.Duration
}

// PolicyConfig stores file-backed policy engine settings.
type PolicyConfig struct {
	FilePath       string
	ReloadInterval time.Duration
}

// SemanticConfig stores semantic inspection settings.
type SemanticConfig struct {
	FailClosed     bool
	MinimumScore   float64
	BlockedIntents []string
}

// AuditConfig stores audit logging settings.
type AuditConfig struct {
	FilePath string
}

// LoadFromEnv reads configuration from environment variables.
func LoadFromEnv() (Config, error) {
	cfg := Config{
		ListenAddr:        getEnv("AEGIS_LISTEN_ADDR", defaultListenAddr),
		UpstreamURL:       os.Getenv("AEGIS_UPSTREAM_URL"),
		UpstreamTimeout:   defaultUpstreamTimeout,
		ReadTimeout:       defaultReadTimeout,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
		WriteTimeout:      defaultWriteTimeout,
		IdleTimeout:       defaultIdleTimeout,
		MaxBodyBytes:      defaultMaxBodyBytes,
		DecisionTimeout:   defaultDecisionTimeout,
		MetricsNamespace:  getEnv("AEGIS_METRICS_NAMESPACE", defaultMetricsNamespace),
		Auth: AuthConfig{
			SharedSecret:     os.Getenv("AEGIS_JWT_SHARED_SECRET"),
			Issuer:           os.Getenv("AEGIS_JWT_ISSUER"),
			Audience:         splitCommaEnv("AEGIS_JWT_AUDIENCE"),
			RequireAudience:  boolEnv("AEGIS_JWT_REQUIRE_AUDIENCE", false),
			AllowedClockSkew: 30 * time.Second,
		},
		Policy: PolicyConfig{
			FilePath:       getEnv("AEGIS_POLICY_FILE", "config/policy.rego"),
			ReloadInterval: defaultPolicyReloadInterval,
		},
		Semantic: SemanticConfig{
			FailClosed:     boolEnv("AEGIS_SEMANTIC_FAIL_CLOSED", defaultSemanticFailClosed),
			MinimumScore:   defaultSemanticMinimumScore,
			BlockedIntents: defaultBlockedIntents(),
		},
		Audit: AuditConfig{
			FilePath: getEnv("AEGIS_AUDIT_FILE", "config/audit.db"),
		},
	}

	var err error
	cfg.UpstreamTimeout, err = durationEnv("AEGIS_UPSTREAM_TIMEOUT", cfg.UpstreamTimeout)
	if err != nil {
		return Config{}, err
	}
	cfg.ReadTimeout, err = durationEnv("AEGIS_READ_TIMEOUT", cfg.ReadTimeout)
	if err != nil {
		return Config{}, err
	}
	cfg.ReadHeaderTimeout, err = durationEnv("AEGIS_READ_HEADER_TIMEOUT", cfg.ReadHeaderTimeout)
	if err != nil {
		return Config{}, err
	}
	cfg.WriteTimeout, err = durationEnv("AEGIS_WRITE_TIMEOUT", cfg.WriteTimeout)
	if err != nil {
		return Config{}, err
	}
	cfg.IdleTimeout, err = durationEnv("AEGIS_IDLE_TIMEOUT", cfg.IdleTimeout)
	if err != nil {
		return Config{}, err
	}
	cfg.DecisionTimeout, err = durationEnv("AEGIS_DECISION_TIMEOUT", cfg.DecisionTimeout)
	if err != nil {
		return Config{}, err
	}
	cfg.Policy.ReloadInterval, err = durationEnv("AEGIS_POLICY_RELOAD_INTERVAL", cfg.Policy.ReloadInterval)
	if err != nil {
		return Config{}, err
	}
	cfg.Auth.AllowedClockSkew, err = durationEnv("AEGIS_JWT_ALLOWED_CLOCK_SKEW", cfg.Auth.AllowedClockSkew)
	if err != nil {
		return Config{}, err
	}
	cfg.MaxBodyBytes, err = int64Env("AEGIS_MAX_BODY_BYTES", cfg.MaxBodyBytes)
	if err != nil {
		return Config{}, err
	}
	cfg.Semantic.MinimumScore, err = float64Env("AEGIS_SEMANTIC_MINIMUM_SCORE", cfg.Semantic.MinimumScore)
	if err != nil {
		return Config{}, err
	}
	if intents := splitCommaEnv("AEGIS_SEMANTIC_BLOCKED_INTENTS"); len(intents) > 0 {
		cfg.Semantic.BlockedIntents = intents
	}

	if cfg.UpstreamURL == "" {
		return Config{}, errors.New("AEGIS_UPSTREAM_URL is required")
	}
	if _, err := url.ParseRequestURI(cfg.UpstreamURL); err != nil {
		return Config{}, fmt.Errorf("invalid AEGIS_UPSTREAM_URL: %w", err)
	}
	if cfg.Auth.SharedSecret == "" {
		return Config{}, errors.New("AEGIS_JWT_SHARED_SECRET is required")
	}
	if cfg.MaxBodyBytes <= 0 {
		return Config{}, errors.New("AEGIS_MAX_BODY_BYTES must be greater than zero")
	}
	if cfg.DecisionTimeout <= 0 {
		return Config{}, errors.New("AEGIS_DECISION_TIMEOUT must be greater than zero")
	}
	if cfg.Semantic.MinimumScore < 0 || cfg.Semantic.MinimumScore > 1 {
		return Config{}, errors.New("AEGIS_SEMANTIC_MINIMUM_SCORE must be between 0 and 1")
	}
	if cfg.Policy.FilePath == "" {
		return Config{}, errors.New("AEGIS_POLICY_FILE is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func splitCommaEnv(key string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func boolEnv(key string, fallback bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func durationEnv(key string, fallback time.Duration) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}
	return parsed, nil
}

func int64Env(key string, fallback int64) (int64, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}
	return parsed, nil
}

func float64Env(key string, fallback float64) (float64, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}
	return parsed, nil
}

func defaultBlockedIntents() []string {
	return []string{
		"prompt_injection",
		"jailbreak_attempt",
		"malicious_code_generation",
		"secret_exfiltration",
	}
}
