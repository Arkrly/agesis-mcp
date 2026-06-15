package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/yourorg/aegis-mcp/internal/config"
)

func TestValidatorAcceptsValidToken(t *testing.T) {
	t.Parallel()

	validator := NewValidator(config.AuthConfig{SharedSecret: "secret"})
	validator.now = func() time.Time { return time.Unix(1_700_000_000, 0).UTC() }

	token := signJWT(t, "secret", map[string]any{
		"sub":      "agent-1",
		"agent_id": "agent-1",
		"roles":    []string{"developer"},
		"exp":      validator.now().Add(time.Hour).Unix(),
	})

	claims, err := validator.Validate(token)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if claims.AgentID != "agent-1" {
		t.Fatalf("AgentID = %q", claims.AgentID)
	}
}

func TestValidatorRejectsExpiredToken(t *testing.T) {
	t.Parallel()

	validator := NewValidator(config.AuthConfig{SharedSecret: "secret"})
	validator.now = func() time.Time { return time.Unix(1_700_000_000, 0).UTC() }

	token := signJWT(t, "secret", map[string]any{
		"sub":      "agent-1",
		"agent_id": "agent-1",
		"exp":      validator.now().Add(-time.Minute).Unix(),
	})

	if _, err := validator.Validate(token); err != ErrExpiredToken {
		t.Fatalf("Validate() error = %v, want %v", err, ErrExpiredToken)
	}
}

func signJWT(t *testing.T, secret string, payload map[string]any) string {
	t.Helper()

	headerBytes, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		t.Fatalf("marshal header: %v", err)
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	headerPart := base64.RawURLEncoding.EncodeToString(headerBytes)
	payloadPart := base64.RawURLEncoding.EncodeToString(payloadBytes)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(headerPart + "." + payloadPart))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return headerPart + "." + payloadPart + "." + signature
}

func FuzzExtractBearerToken(f *testing.F) {
	f.Add("Bearer ")
	f.Add("Bearer abc.def.ghi")
	f.Add("bearer token123")
	f.Add("Basic xyz")
	f.Add("")
	f.Add("  Bearer   token  ")

	f.Fuzz(func(t *testing.T, header string) {
		token, err := ExtractBearerToken(header)
		if err == nil {
			if token == "" {
				t.Errorf("ExtractBearerToken returned empty token with no error for header %q", header)
			}
		}
	})
}
