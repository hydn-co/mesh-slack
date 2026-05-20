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
	Profile           slackUserProfile `json:"profile"`
	ID                string           `json:"id"`
	Name              string           `json:"name"`
	RealName          string           `json:"real_name"`
	Deleted           bool             `json:"deleted"`
	IsBot             bool             `json:"is_bot"`
	IsAppUser         bool             `json:"is_app_user"`
	IsRestricted      bool             `json:"is_restricted"`
	IsUltraRestricted bool             `json:"is_ultra_restricted"`
}

type ListUsersResult struct {
	slackapi.ResponseEnvelope
	ResponseMetadata responseMetadata `json:"response_metadata"`
	Members          []SlackUser      `json:"members"`
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
