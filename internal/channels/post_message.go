package channels

import (
	"context"

	"github.com/hydn-co/mesh-slack/internal/endpoints"
	slackapi "github.com/hydn-co/mesh-slack/internal/slack_api"
)

// PostMessage sends a plain text message via Slack's chat.postMessage API.
func PostMessage(ctx context.Context, token, channelID, text string) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	payload := map[string]string{
		"channel": channelID,
		"text":    text,
	}

	req, err := slackapi.NewJSONRequest(ctx, endpoints.SlackChatPostMessage, token, payload)
	if err != nil {
		return err
	}

	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	return slackapi.Do(req)
}
