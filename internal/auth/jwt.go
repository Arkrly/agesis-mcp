package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/yourorg/aegis-mcp/internal/config"
)

var (
	ErrMissingAuthorization = errors.New("missing authorization header")
	ErrInvalidAuthorization = errors.New("invalid authorization header")
	ErrInvalidToken         = errors.New("invalid jwt token")
	ErrExpiredToken         = errors.New("jwt token expired")
	ErrTokenNotYetValid     = errors.New("jwt token not yet valid")
	ErrMissingAgentID       = errors.New("jwt token missing agent_id")
)

type contextKey string

const claimsContextKey contextKey = "auth.claims"

// Claims is the authorization context extracted from the JWT.
type Claims struct {
	Subject   string
	AgentID   string
	Roles     []string
	Issuer    string
	Audience  []string
	ExpiresAt *time.Time
	NotBefore *time.Time
	IssuedAt  *time.Time
	Metadata  map[string]any
}

// Validator validates HS256 JWT tokens.
type Validator struct {
	secret           []byte
	issuer           string
	audience         []string
	requireAudience  bool
	allowedClockSkew time.Duration
	now              func() time.Time
}

func NewValidator(cfg config.AuthConfig) *Validator {
	return &Validator{
		secret:           []byte(cfg.SharedSecret),
		issuer:           cfg.Issuer,
		audience:         append([]string(nil), cfg.Audience...),
		requireAudience:  cfg.RequireAudience,
		allowedClockSkew: cfg.AllowedClockSkew,
		now:              time.Now,
	}
}

func (v *Validator) ClaimsFromRequest(r *http.Request) (Claims, error) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	token, err := ExtractBearerToken(authHeader)
	if err != nil {
		return Claims{}, err
	}
	return v.Validate(token)
}

func ExtractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", ErrMissingAuthorization
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		return "", ErrInvalidAuthorization
	}
	return parts[1], nil
}

func (v *Validator) Validate(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, ErrInvalidToken
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Claims{}, fmt.Errorf("%w: decode header", ErrInvalidToken)
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, fmt.Errorf("%w: decode payload", ErrInvalidToken)
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return Claims{}, fmt.Errorf("%w: decode signature", ErrInvalidToken)
	}

	var header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return Claims{}, fmt.Errorf("%w: parse header", ErrInvalidToken)
	}
	if header.Alg != "HS256" {
		return Claims{}, fmt.Errorf("%w: unsupported alg %q", ErrInvalidToken, header.Alg)
	}

	mac := hmac.New(sha256.New, v.secret)
	mac.Write([]byte(parts[0] + "." + parts[1]))
	expected := mac.Sum(nil)
	if !hmac.Equal(signature, expected) {
		return Claims{}, ErrInvalidToken
	}

	var payload map[string]any
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return Claims{}, fmt.Errorf("%w: parse payload", ErrInvalidToken)
	}

	claims, err := v.mapClaims(payload)
	if err != nil {
		return Claims{}, err
	}
	if err := v.validateRegisteredClaims(claims); err != nil {
		return Claims{}, err
	}
	return claims, nil
}

func ContextWithClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(Claims)
	return claims, ok
}

func (v *Validator) mapClaims(payload map[string]any) (Claims, error) {
	claims := Claims{
		Metadata: make(map[string]any),
	}
	if sub, _ := payload["sub"].(string); sub != "" {
		claims.Subject = sub
	}
	if agentID, _ := payload["agent_id"].(string); agentID != "" {
		claims.AgentID = agentID
	} else if claims.Subject != "" {
		claims.AgentID = claims.Subject
	}
	if claims.AgentID == "" {
		return Claims{}, ErrMissingAgentID
	}
	if issuer, _ := payload["iss"].(string); issuer != "" {
		claims.Issuer = issuer
	}
	claims.Roles = parseStringList(payload["roles"])
	claims.Audience = parseAudience(payload["aud"])
	if ts, ok := parseUnixTime(payload["exp"]); ok {
		claims.ExpiresAt = &ts
	}
	if ts, ok := parseUnixTime(payload["nbf"]); ok {
		claims.NotBefore = &ts
	}
	if ts, ok := parseUnixTime(payload["iat"]); ok {
		claims.IssuedAt = &ts
	}
	if metadata, ok := payload["metadata"].(map[string]any); ok {
		claims.Metadata = metadata
	}
	return claims, nil
}

func (v *Validator) validateRegisteredClaims(claims Claims) error {
	now := v.now()
	if claims.ExpiresAt != nil && now.After(claims.ExpiresAt.Add(v.allowedClockSkew)) {
		return ErrExpiredToken
	}
	if claims.NotBefore != nil && now.Add(v.allowedClockSkew).Before(*claims.NotBefore) {
		return ErrTokenNotYetValid
	}
	if v.issuer != "" && claims.Issuer != v.issuer {
		return fmt.Errorf("%w: issuer mismatch", ErrInvalidToken)
	}
	if v.requireAudience {
		if len(v.audience) == 0 {
			return errors.New("audience required but not configured")
		}
		if !intersects(claims.Audience, v.audience) {
			return fmt.Errorf("%w: audience mismatch", ErrInvalidToken)
		}
	}
	return nil
}

func parseStringList(value any) []string {
	switch typed := value.(type) {
	case []any:
		result := make([]string, 0, len(typed))
		for _, item := range typed {
			if text, ok := item.(string); ok && text != "" {
				result = append(result, text)
			}
		}
		return result
	case string:
		if typed == "" {
			return nil
		}
		return []string{typed}
	default:
		return nil
	}
}

func parseAudience(value any) []string {
	return parseStringList(value)
}

func parseUnixTime(value any) (time.Time, bool) {
	switch typed := value.(type) {
	case float64:
		return time.Unix(int64(typed), 0).UTC(), true
	case int64:
		return time.Unix(typed, 0).UTC(), true
	default:
		return time.Time{}, false
	}
}

func intersects(left, right []string) bool {
	for _, a := range left {
		for _, b := range right {
			if a == b {
				return true
			}
		}
	}
	return false
}
