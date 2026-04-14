package slackapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type responseEnvelope struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

// EnsureContextActive returns early when the provided context has been canceled.
func EnsureContextActive(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("operation canceled: %w", err)
	}

	return nil
}

func NewFormRequest(ctx context.Context, endpoint, token string, data url.Values) (*http.Request, error) {
	if err := EnsureContextActive(ctx); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func NewJSONRequest(ctx context.Context, endpoint, token string, payload any) (*http.Request, error) {
	if err := EnsureContextActive(ctx); err != nil {
		return nil, err
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func Do(req *http.Request) error {
	if err := EnsureContextActive(req.Context()); err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if cerr := req.Context().Err(); cerr != nil {
			return fmt.Errorf("operation canceled: %w", cerr)
		}
		return fmt.Errorf("API request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if err := EnsureContextActive(req.Context()); err != nil {
		return err
	}

	var result responseEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to parse API response: %w", err)
	}

	if !result.OK {
		if result.Error == "" {
			return fmt.Errorf("slack API call failed")
		}
		return fmt.Errorf("slack API error: %s", result.Error)
	}

	return nil
}
