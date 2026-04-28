package users

import (
	"context"
	"net/url"
	"strconv"

	"github.com/hydn-co/mesh-slack/internal/endpoints"
	slackapi "github.com/hydn-co/mesh-slack/internal/slack_api"
)

type SlackUserConversation struct {
	ID string `json:"id"`
}

type ListUserConversationsResult struct {
	slackapi.ResponseEnvelope
	Channels         []SlackUserConversation `json:"channels"`
	ResponseMetadata responseMetadata        `json:"response_metadata"`
}

// ListUserConversations lists channels for a specific Slack user.
func ListUserConversations(ctx context.Context, token, userID, cursor string) (*ListUserConversationsResult, error) {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return nil, err
	}

	data := url.Values{
		"exclude_archived": {"false"},
		"limit":            {strconv.Itoa(slackapi.MaxPageLimit)},
		"types":            {"public_channel,private_channel"},
		"user":             {userID},
	}
	if cursor != "" {
		data.Set("cursor", cursor)
	}

	req, err := slackapi.NewFormRequest(ctx, endpoints.SlackUsersConversations, token, data)
	if err != nil {
		return nil, err
	}

	var response ListUserConversationsResult
	if err := slackapi.Do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
