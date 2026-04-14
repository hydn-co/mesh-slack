package channels

import (
	"context"
	"net/url"

	"github.com/hydn-co/mesh-slack/internal/endpoints"
	slackapi "github.com/hydn-co/mesh-slack/internal/slack_api"
)

// ValidateExists verifies the channel exists and is accessible with the bot token.
func ValidateExists(ctx context.Context, token, channelID string) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	data := url.Values{"channel": {channelID}}

	req, err := slackapi.NewFormRequest(ctx, endpoints.SlackConversationsInfo, token, data)
	if err != nil {
		return err
	}

	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	return slackapi.Do(req)
}
