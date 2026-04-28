package slackapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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

const maxRateLimitRetries = 5

const MaxPageLimit = 999

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

	for attempt := 0; attempt <= maxRateLimitRetries; attempt++ {
		resp, err := doRequestAttempt(req, attempt)
		if err != nil {
			if cerr := req.Context().Err(); cerr != nil {
				return fmt.Errorf("operation canceled: %w", cerr)
			}
			return fmt.Errorf("API request failed: %w", err)
		}

		body, readErr := readResponseBody(req.Context(), resp)
		if readErr != nil {
			return readErr
		}

		if err := EnsureContextActive(req.Context()); err != nil {
			return err
		}

		var envelope ResponseEnvelope
		envelopeErr := json.Unmarshal(body, &envelope)

		if resp.StatusCode == http.StatusTooManyRequests {
			if attempt == maxRateLimitRetries {
				retryAfter := parseRetryAfterSeconds(resp.Header.Get("Retry-After"))
				if retryAfter > 0 {
					return fmt.Errorf("slack API %s was rate limited after %d retries; retry after %d seconds", slackMethodName(req), maxRateLimitRetries, retryAfter)
				}
				return fmt.Errorf("slack API %s was rate limited after %d retries", slackMethodName(req), maxRateLimitRetries)
			}

			if err := waitForRetry(req.Context(), rateLimitDelay(resp.Header.Get("Retry-After"), attempt)); err != nil {
				return err
			}

			continue
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

	return fmt.Errorf("slack API %s failed unexpectedly", slackMethodName(req))
}

func doRequestAttempt(req *http.Request, attempt int) (*http.Response, error) {
	if attempt > 0 {
		if req.Body != nil && req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				return nil, fmt.Errorf("failed to reset request body for retry: %w", err)
			}

			req.Body = body
		}
	}

	return http.DefaultClient.Do(req)
}

func readResponseBody(ctx context.Context, resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, fmt.Errorf("response is nil")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.WarnContext(ctx, "failed to close response body", slog.Any("error", err))
		}
	}()

	if err := EnsureContextActive(ctx); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %w", err)
	}

	return body, nil
}

func waitForRetry(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		delay = time.Second
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return fmt.Errorf("operation canceled: %w", ctx.Err())
	case <-timer.C:
		return nil
	}
}

func rateLimitDelay(value string, attempt int) time.Duration {
	if retryAfter := parseRetryAfterSeconds(value); retryAfter > 0 {
		return time.Duration(retryAfter) * time.Second
	}

	return time.Duration(1<<attempt) * time.Second
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

	if description, ok := slackErrorDescriptions[code]; ok {
		return description
	}

	return code
}
