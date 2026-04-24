package users

import (
	"context"
	"fmt"
	"strings"

	slackapi "github.com/hydn-co/mesh-slack/internal/slack_api"
)

// FindUserIDByEmail resolves a Slack user ID by scanning workspace members for a
// matching email address. Returns an error if no match is found.
func FindUserIDByEmail(ctx context.Context, token, email string) (string, error) {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return "", err
	}

	normalized := strings.TrimSpace(email)
	if normalized == "" {
		return "", fmt.Errorf("email is required")
	}

	cursor := ""
	for {
		result, err := ListUsers(ctx, token, cursor)
		if err != nil {
			return "", fmt.Errorf("failed to list Slack users: %w", err)
		}

		for _, user := range result.Members {
			if strings.EqualFold(strings.TrimSpace(user.Profile.Email), normalized) {
				if user.ID == "" {
					return "", fmt.Errorf("resolved Slack user for email %s has no user ID", normalized)
				}
				return user.ID, nil
			}
		}

		cursor = strings.TrimSpace(result.ResponseMetadata.NextCursor)
		if cursor == "" {
			break
		}
	}

	return "", fmt.Errorf("slack user not found for email %s", normalized)
}

// ResolveUserIDsByEmails resolves multiple Slack user IDs in a single paginated
// scan of the workspace member list. Returns an error if any email is not found.
func ResolveUserIDsByEmails(ctx context.Context, token string, emails []string) ([]string, error) {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return nil, err
	}

	// Build a lookup map: normalised email → original (for error messages).
	want := make(map[string]string, len(emails))
	for _, e := range emails {
		n := strings.ToLower(strings.TrimSpace(e))
		if n == "" {
			return nil, fmt.Errorf("email list contains an empty entry")
		}
		want[n] = e
	}

	resolved := make(map[string]string, len(emails)) // normalised email → user ID

	cursor := ""
	for {
		result, err := ListUsers(ctx, token, cursor)
		if err != nil {
			return nil, fmt.Errorf("failed to list Slack users: %w", err)
		}

		for _, user := range result.Members {
			n := strings.ToLower(strings.TrimSpace(user.Profile.Email))
			if _, needed := want[n]; needed {
				if user.ID == "" {
					return nil, fmt.Errorf("resolved Slack user for email %s has no user ID", want[n])
				}
				resolved[n] = user.ID
			}
		}

		if len(resolved) == len(want) {
			break // all found — no need to continue paginating
		}

		cursor = strings.TrimSpace(result.ResponseMetadata.NextCursor)
		if cursor == "" {
			break
		}
	}

	// Report any emails that had no match.
	if len(resolved) != len(want) {
		var missing []string
		for n, orig := range want {
			if _, ok := resolved[n]; !ok {
				missing = append(missing, orig)
			}
		}
		return nil, fmt.Errorf("slack users not found for emails: %s", strings.Join(missing, ", "))
	}

	// Return IDs in the same order as the input emails.
	ids := make([]string, 0, len(emails))
	for _, e := range emails {
		ids = append(ids, resolved[strings.ToLower(strings.TrimSpace(e))])
	}
	return ids, nil
}
