package slackapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var slackErrorDescriptions = map[string]string{
	"channel_not_found": "channel was not found or is not accessible to the bot",
	"not_in_channel":    "bot is not a member of the channel",
	"is_archived":       "channel is archived",
	"missing_scope":     "bot token is missing a required Slack scope",
	"invalid_auth":      "Slack authentication failed; verify the bot token",
	"not_authed":        "Slack authentication failed; no bot token was provided",
	"token_revoked":     "Slack authentication failed; the bot token has been revoked",
	"account_inactive":  "Slack authentication failed; the workspace account is inactive",
	"msg_too_long":      "message exceeds Slack's maximum message length",
	"rate_limited":      "Slack rate limited the request",
}

type ResponseEnvelope struct {
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

func Do(req *http.Request, response any) error {
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
		err = resp.Body.Close()
		if err != nil {
			slog.WarnContext(req.Context(), "failed to close response body", slog.Any("error", err))
		}
	}()

	if err := EnsureContextActive(req.Context()); err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read API response: %w", err)
	}

	if err := EnsureContextActive(req.Context()); err != nil {
		return err
	}

	var envelope ResponseEnvelope
	envelopeErr := json.Unmarshal(body, &envelope)

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := parseRetryAfterSeconds(resp.Header.Get("Retry-After"))
		if retryAfter > 0 {
			return fmt.Errorf("slack API %s was rate limited; retry after %d seconds", slackMethodName(req), retryAfter)
		}
		return fmt.Errorf("slack API %s was rate limited", slackMethodName(req))
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if envelopeErr == nil && envelope.Error != "" {
			return fmt.Errorf("slack API %s failed with status %d: %s", slackMethodName(req), resp.StatusCode, describeSlackError(envelope.Error))
		}

		trimmedBody := strings.TrimSpace(string(body))
		if trimmedBody != "" {
			return fmt.Errorf("slack API %s failed with status %d: %s", slackMethodName(req), resp.StatusCode, trimmedBody)
		}

		return fmt.Errorf("slack API %s failed with status %d", slackMethodName(req), resp.StatusCode)
	}

	if envelopeErr != nil {
		return fmt.Errorf("failed to parse API response: %w", envelopeErr)
	}

	if !envelope.OK {
		if envelope.Error == "" {
			return fmt.Errorf("slack API %s failed", slackMethodName(req))
		}
		return fmt.Errorf("slack API %s failed: %s", slackMethodName(req), describeSlackError(envelope.Error))
	}

	if response != nil {
		if err := json.Unmarshal(body, response); err != nil {
			return fmt.Errorf("failed to parse typed API response: %w", err)
		}
	}

	return nil
}

func parseRetryAfterSeconds(value string) int {
	if value == "" {
		return 0
	}

	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		return 0
	}

	return seconds
}

func slackMethodName(req *http.Request) string {
	if req == nil || req.URL == nil {
		return "request"
	}

	trimmed := strings.Trim(req.URL.Path, "/")
	if trimmed == "" {
		return req.URL.String()
	}

	parts := strings.Split(trimmed, "/")
	return parts[len(parts)-1]
}

func describeSlackError(code string) string {
	if code == "" {
		return "unknown Slack error"
	}

	descriptions := maps.Clone(slackErrorDescriptions)
	if description, ok := descriptions[code]; ok {
		return description
	}

	return code
}
