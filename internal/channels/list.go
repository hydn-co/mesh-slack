package channels

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

type slackChannelPurpose struct {
	Value string `json:"value"`
}

type SlackChannel struct {
	ID         string              `json:"id"`
	Name       string              `json:"name"`
	IsArchived bool                `json:"is_archived"`
	IsPrivate  bool                `json:"is_private"`
	Purpose    slackChannelPurpose `json:"purpose"`
}

type ListChannelsResult struct {
	slackapi.ResponseEnvelope
	Channels         []SlackChannel   `json:"channels"`
	ResponseMetadata responseMetadata `json:"response_metadata"`
}

// ListChannels lists Slack channels visible to the bot token.
func ListChannels(ctx context.Context, token, cursor string) (*ListChannelsResult, error) {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return nil, err
	}

	data := url.Values{
		"exclude_archived": {"false"},
		"limit":            {strconv.Itoa(slackapi.MaxPageLimit)},
		"types":            {"public_channel,private_channel"},
	}
	if cursor != "" {
		data.Set("cursor", cursor)
	}

	req, err := slackapi.NewFormRequest(ctx, endpoints.SlackConversationsList, token, data)
	if err != nil {
		return nil, err
	}

	var response ListChannelsResult
	if err := slackapi.Do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
