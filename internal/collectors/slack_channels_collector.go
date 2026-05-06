package collectors

import (
	"context"
	"fmt"

	"github.com/fgrzl/enumerators"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/entities"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/types"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
	"github.com/hydn-co/mesh-slack/internal/channels"
	"github.com/hydn-co/mesh-slack/internal/options"
	slackapi "github.com/hydn-co/mesh-slack/internal/slack_api"
)

// SlackChannelsCollector collects channels from Slack workspaces and emits them
// as catalog entities.
type SlackChannelsCollector struct {
	*connector.TypedFeatureContext[*options.SlackChannelsCollectorOptions, *connector.NoPayload]
	token string
	state connectorutil.FeatureState
}

// NewSlackChannelsCollector constructs a SlackChannelsCollector.
func NewSlackChannelsCollector(ctx *connector.TypedFeatureContext[*options.SlackChannelsCollectorOptions, *connector.NoPayload]) runner.Feature {
	return &SlackChannelsCollector{TypedFeatureContext: ctx}
}

// Init prepares the collector for operation.
func (c *SlackChannelsCollector) Init(ctx context.Context) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	token, err := connectorutil.ExtractAPIKey(c.GetCredentials())
	if err != nil {
		return fmt.Errorf("failed to extract bot token: %w", err)
	}

	c.token = token
	c.state.MarkReady()
	return nil
}

// Start begins collecting channels from the Slack workspace.
func (c *SlackChannelsCollector) Start(ctx context.Context) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	if err := c.state.RequireReady(); err != nil {
		return err
	}

	channelEnum := channels.ChannelEnumerator(ctx, c.token)
	if err := enumerators.ForEach(channelEnum, func(channel channels.SlackChannel) error {
		if err := slackapi.EnsureContextActive(ctx); err != nil {
			return err
		}

		entity := &entities.Channel{
			Metadata:    types.EntityMetadata{Space: spaces.Channels},
			ChannelRef:  channel.ID,
			Name:        channel.Name,
			Description: channel.Purpose.Value,
			Archived:    channel.IsArchived,
			Private:     channel.IsPrivate,
		}

		if err := c.Emit(ctx, entity); err != nil {
			return fmt.Errorf("failed to emit channel %s: %w", channel.ID, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to enumerate channels: %w", err)
	}

	return nil
}

// Stop halts channel collection and releases resources.
func (c *SlackChannelsCollector) Stop(ctx context.Context) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	if err := c.state.RequireReady(); err != nil {
		return err
	}

	c.state.Reset()
	c.token = ""
	return nil
}
