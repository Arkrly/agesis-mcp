package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	secret := "test-secret-key-32-chars-long-12345"
	
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	
	payload := map[string]any{
		"sub":      "agent-007",
		"agent_id": "agent-007",
		"roles":    []string{"developer"},
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	
	headerJSON, _ := json.Marshal(header)
	payloadJSON, _ := json.Marshal(payload)
	
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)
	
	unsignedToken := headerEncoded + "." + payloadEncoded
	
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(unsignedToken))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	
	fmt.Println(unsignedToken + "." + signature)
}
