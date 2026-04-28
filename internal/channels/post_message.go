package channels

import (
	"context"

	"github.com/hydn-co/mesh-slack/internal/endpoints"
	slackapi "github.com/hydn-co/mesh-slack/internal/slack_api"
)

type postMessageRequest struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type postMessageMessage struct {
	Text     string `json:"text"`
	Type     string `json:"type"`
	BotID    string `json:"bot_id"`
	ThreadTS string `json:"thread_ts,omitempty"`
}

// PostMessageResult captures the useful fields returned by Slack chat.postMessage.
type PostMessageResult struct {
	slackapi.ResponseEnvelope
	Channel string             `json:"channel"`
	TS      string             `json:"ts"`
	Message postMessageMessage `json:"message"`
}

// PostMessage sends a plain text message via Slack's chat.postMessage API.
func PostMessage(ctx context.Context, token, channelID, text string) (*PostMessageResult, error) {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return nil, err
	}

	payload := postMessageRequest{
		Channel: channelID,
		Text:    text,
	}

	req, err := slackapi.NewJSONRequest(ctx, endpoints.SlackChatPostMessage, token, payload)
	if err != nil {
		return nil, err
	}

	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return nil, err
	}

	var response PostMessageResult
	if err := slackapi.Do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
