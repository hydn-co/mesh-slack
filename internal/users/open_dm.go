package users

import (
	"context"
	"fmt"
	"strings"

	"github.com/hydn-co/mesh-slack/internal/endpoints"
	slackapi "github.com/hydn-co/mesh-slack/internal/slack_api"
)

type openDMRequest struct {
	Users string `json:"users"`
}

type openDMChannel struct {
	ID string `json:"id"`
}

type openDMResponse struct {
	slackapi.ResponseEnvelope
	Channel openDMChannel `json:"channel"`
}

// OpenDM opens (or resolves) a direct message or group DM channel for the given
// Slack user IDs and returns the channel ID for use with channels.PostMessage.
// One ID opens a 1:1 DM; two to eight IDs open a group DM (MPIM).
func OpenDM(ctx context.Context, token string, userIDs []string) (string, error) {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return "", err
	}

	body := openDMRequest{Users: strings.Join(userIDs, ",")}

	req, err := slackapi.NewJSONRequest(ctx, endpoints.SlackConversationsOpen, token, body)
	if err != nil {
		return "", err
	}

	var response openDMResponse
	if err := slackapi.Do(req, &response); err != nil {
		return "", err
	}

	if response.Channel.ID == "" {
		return "", fmt.Errorf("conversations.open returned empty channel ID for users %s", strings.Join(userIDs, ", "))
	}

	return response.Channel.ID, nil
}
