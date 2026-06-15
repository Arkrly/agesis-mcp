package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// UpstreamClient forwards MCP requests to the configured backend.
type UpstreamClient interface {
	Forward(ctx context.Context, body []byte, authHeader string, sessionID string) (*http.Response, error)
}

type HTTPUpstreamClient struct {
	targetURL string
	client    *http.Client
}

func NewHTTPUpstreamClient(targetURL string, timeout time.Duration) *HTTPUpstreamClient {
	return &HTTPUpstreamClient{
		targetURL: targetURL,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *HTTPUpstreamClient) Forward(ctx context.Context, body []byte, authHeader string, sessionID string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.targetURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create upstream request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	if sessionID != "" {
		req.Header.Set("Mcp-Session-Id", sessionID)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("forward request: %w", err)
	}

	return resp, nil
}

func CopyResponse(w http.ResponseWriter, resp *http.Response) error {
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("copy response body: %w", err)
	}

	return nil
}
