package users

import (
	"context"
	"net/url"
	"strconv"

	"github.com/hydn-co/mesh-slack/internal/endpoints"
	slackapi "github.com/hydn-co/mesh-slack/internal/slack_api"
)

type responseMetadata struct {
	NextCursor string `json:"next_cursor"`
}

type slackUserProfile struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	RealName    string `json:"real_name"`
	Title       string `json:"title"`
}

type SlackUser struct {
	ID                string           `json:"id"`
	Name              string           `json:"name"`
	Deleted           bool             `json:"deleted"`
	IsBot             bool             `json:"is_bot"`
	IsAppUser         bool             `json:"is_app_user"`
	IsRestricted      bool             `json:"is_restricted"`
	IsUltraRestricted bool             `json:"is_ultra_restricted"`
	RealName          string           `json:"real_name"`
	Profile           slackUserProfile `json:"profile"`
}

type ListUsersResult struct {
	slackapi.ResponseEnvelope
	Members          []SlackUser      `json:"members"`
	ResponseMetadata responseMetadata `json:"response_metadata"`
}

// ListUsers lists Slack workspace users visible to the bot token.
func ListUsers(ctx context.Context, token, cursor string) (*ListUsersResult, error) {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return nil, err
	}

	data := url.Values{
		"limit": {strconv.Itoa(slackapi.MaxPageLimit)},
	}
	if cursor != "" {
		data.Set("cursor", cursor)
	}

	req, err := slackapi.NewFormRequest(ctx, endpoints.SlackUsersList, token, data)
	if err != nil {
		return nil, err
	}

	var response ListUsersResult
	if err := slackapi.Do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
