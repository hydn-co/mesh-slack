package collectors

import (
	"context"

	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
	"github.com/hydn-co/mesh-slack/internal/helpers"
	"github.com/hydn-co/mesh-slack/internal/options"
	slackapi "github.com/hydn-co/mesh-slack/internal/slack_api"
)

// SlackChannelsCollector collects channels from Slack workspaces and emits them
// as catalog entities.
type SlackChannelsCollector struct {
	*connector.TypedFeatureContext[*options.SlackChannelsCollectorOptions]
	initialized bool
}

// NewSlackChannelsCollector constructs a SlackChannelsCollector.
func NewSlackChannelsCollector(ctx *connector.TypedFeatureContext[*options.SlackChannelsCollectorOptions]) runner.Feature {
	return &SlackChannelsCollector{TypedFeatureContext: ctx}
}

// Init prepares the collector for operation.
func (c *SlackChannelsCollector) Init(ctx context.Context) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	c.initialized = true
	return nil
}

// Start begins collecting channels from the Slack workspace.
func (c *SlackChannelsCollector) Start(ctx context.Context) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	if err := helpers.CheckInitialized(c.initialized); err != nil {
		return err
	}

	return nil
}

// Stop halts channel collection and releases resources.
func (c *SlackChannelsCollector) Stop(ctx context.Context) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	if err := helpers.CheckInitialized(c.initialized); err != nil {
		return err
	}

	c.initialized = false
	return nil
}
