package channels

import (
	"context"
	"net/url"

	"github.com/hydn-co/mesh-slack/internal/endpoints"
	slackapi "github.com/hydn-co/mesh-slack/internal/slack_api"
)

type validateExistsRequest struct {
	ChannelID string
}

func (r validateExistsRequest) formValues() url.Values {
	return url.Values{"channel": {r.ChannelID}}
}

type validateExistsChannel struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	IsChannel  bool   `json:"is_channel"`
	IsGroup    bool   `json:"is_group"`
	IsPrivate  bool   `json:"is_private"`
	IsArchived bool   `json:"is_archived"`
}

type validateExistsResponse struct {
	slackapi.ResponseEnvelope
	Channel validateExistsChannel `json:"channel"`
}

// ValidateExists verifies the channel exists and is accessible with the bot token.
func ValidateExists(ctx context.Context, token, channelID string) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	request := validateExistsRequest{ChannelID: channelID}

	req, err := slackapi.NewFormRequest(ctx, endpoints.SlackConversationsInfo, token, request.formValues())
	if err != nil {
		return err
	}

	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	var response validateExistsResponse
	return slackapi.Do(req, &response)
}
