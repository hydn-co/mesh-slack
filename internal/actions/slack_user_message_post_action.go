package actions

import (
	"context"
	"fmt"

	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
	"github.com/hydn-co/mesh-slack/internal/channels"
	"github.com/hydn-co/mesh-slack/internal/options"
	"github.com/hydn-co/mesh-slack/internal/payloads"
	slackapi "github.com/hydn-co/mesh-slack/internal/slack_api"
	"github.com/hydn-co/mesh-slack/internal/users"
)

// SlackUserMessagePostAction posts messages to a Slack user via direct message.
type SlackUserMessagePostAction struct {
	*connector.TypedFeatureContext[*options.SlackUserMessagePostActionOptions, *payloads.SlackChannelMessagePostPayload]
	token       string
	dmChannelID string
	message     string
	state       connectorutil.FeatureState
}

// NewSlackUserMessagePostAction constructs a SlackUserMessagePostAction.
func NewSlackUserMessagePostAction(ctx *connector.TypedFeatureContext[*options.SlackUserMessagePostActionOptions, *payloads.SlackChannelMessagePostPayload]) runner.Feature {
	return &SlackUserMessagePostAction{TypedFeatureContext: ctx}
}

// Init prepares the action by validating credentials, resolving the recipient(s), and verifying the message.
func (p *SlackUserMessagePostAction) Init(ctx context.Context) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	opts := p.GetOptions()
	if opts == nil || len(opts.Emails) == 0 {
		return fmt.Errorf("at least one email is required")
	}
	if len(opts.Emails) > 8 {
		return fmt.Errorf("at most 8 recipients are supported, got %d", len(opts.Emails))
	}

	payload := p.GetPayload()
	if payload == nil {
		return fmt.Errorf("message payload is required")
	}
	message, err := verifyMessage(payload.Message)
	if err != nil {
		return err
	}

	token, err := connectorutil.ExtractAPIKey(p.GetCredentials())
	if err != nil {
		return fmt.Errorf("failed to extract bot token: %w", err)
	}

	userIDs, err := users.ResolveUserIDsByEmails(ctx, token, opts.Emails)
	if err != nil {
		return err
	}

	dmChannelID, err := users.OpenDM(ctx, token, userIDs)
	if err != nil {
		return fmt.Errorf("failed to open conversation: %w", err)
	}

	p.token = token
	p.dmChannelID = dmChannelID
	p.message = message
	p.state.MarkReady()

	return nil
}

// Start posts the message to the resolved DM channel.
func (p *SlackUserMessagePostAction) Start(ctx context.Context) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	if err := p.state.RequireReady(); err != nil {
		return err
	}

	_, err := channels.PostMessage(ctx, p.token, p.dmChannelID, p.message)
	return err
}

// Stop halts event processing and releases resources.
func (p *SlackUserMessagePostAction) Stop(ctx context.Context) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	if err := p.state.RequireReady(); err != nil {
		return err
	}

	p.state.Reset()
	p.token = ""
	p.dmChannelID = ""
	p.message = ""

	return nil
}
