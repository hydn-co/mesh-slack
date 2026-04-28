package actions

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
	"github.com/hydn-co/mesh-slack/internal/channels"
	"github.com/hydn-co/mesh-slack/internal/credentials"
	"github.com/hydn-co/mesh-slack/internal/helpers"
	"github.com/hydn-co/mesh-slack/internal/options"
	"github.com/hydn-co/mesh-slack/internal/payloads"
	slackapi "github.com/hydn-co/mesh-slack/internal/slack_api"
)

// SlackChannelMessagePostAction posts messages to Slack channels.
type SlackChannelMessagePostAction struct {
	*connector.TypedFeatureContext[*options.SlackChannelMessagePostActionOptions, *payloads.SlackChannelMessagePostPayload]
	token       string
	message     string
	initialized bool
}

// NewSlackChannelMessagePostAction constructs a SlackChannelMessagePostAction.
func NewSlackChannelMessagePostAction(ctx *connector.TypedFeatureContext[*options.SlackChannelMessagePostActionOptions, *payloads.SlackChannelMessagePostPayload]) runner.Feature {
	return &SlackChannelMessagePostAction{TypedFeatureContext: ctx}
}

// Init prepares the action for operation by validating credentials, channel, and message.
func (p *SlackChannelMessagePostAction) Init(ctx context.Context) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	opts := p.GetOptions()
	if opts == nil || opts.ChannelID == "" {
		return fmt.Errorf("channel_id is required")
	}

	payload := p.GetPayload()
	if payload == nil {
		return fmt.Errorf("message payload is required")
	}
	message, err := verifyMessage(payload.Message)
	if err != nil {
		return err
	}

	token, err := credentials.ExtractToken(p.GetCredentials())
	if err != nil {
		return fmt.Errorf("failed to extract bot token: %w", err)
	}

	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	if err := channels.ValidateExists(ctx, token, opts.ChannelID); err != nil {
		return fmt.Errorf("channel validation failed: %w", err)
	}

	p.token = token
	p.message = message
	p.initialized = true

	return nil
}

// Start begins processing catalog events and posting messages to Slack.
func (p *SlackChannelMessagePostAction) Start(ctx context.Context) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	if err := helpers.CheckInitialized(p.initialized); err != nil {
		return err
	}

	opts := p.GetOptions()
	if opts == nil || opts.ChannelID == "" {
		return fmt.Errorf("channel_id is required")
	}

	_, err := channels.PostMessage(ctx, p.token, opts.ChannelID, p.message)
	return err
}

// Stop halts event processing and releases resources.
func (p *SlackChannelMessagePostAction) Stop(ctx context.Context) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	if err := helpers.CheckInitialized(p.initialized); err != nil {
		return err
	}

	p.initialized = false
	p.token = ""
	p.message = ""

	return nil
}

// ============= Private helpers below =============

func verifyMessage(text string) (string, error) {
	if text == "" {
		return "", fmt.Errorf("message cannot be empty")
	}

	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return "", fmt.Errorf("message cannot be empty")
	}

	charCount := utf8.RuneCountInString(trimmed)
	if charCount > 4000 {
		return "", fmt.Errorf("message exceeds maximum length of 4000 characters (got %d)", charCount)
	}

	if !utf8.ValidString(trimmed) {
		return "", fmt.Errorf("message contains invalid UTF-8")
	}

	return trimmed, nil
}
