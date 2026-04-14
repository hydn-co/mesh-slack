package channels

import (
	"context"
	"net/url"
	"strconv"

	"github.com/hydn-co/mesh-slack/internal/endpoints"
	slackapi "github.com/hydn-co/mesh-slack/internal/slack_api"
)

type ListMembersResult struct {
	slackapi.ResponseEnvelope
	Members          []string         `json:"members"`
	ResponseMetadata responseMetadata `json:"response_metadata"`
}

// ListMembers lists the member account IDs for a Slack channel.
func ListMembers(ctx context.Context, token, channelID, cursor string) (*ListMembersResult, error) {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return nil, err
	}

	data := url.Values{
		"channel": {channelID},
		"limit":   {strconv.Itoa(slackapi.MaxPageLimit)},
	}
	if cursor != "" {
		data.Set("cursor", cursor)
	}

	req, err := slackapi.NewFormRequest(ctx, endpoints.SlackConversationsMembers, token, data)
	if err != nil {
		return nil, err
	}

	var response ListMembersResult
	if err := slackapi.Do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
