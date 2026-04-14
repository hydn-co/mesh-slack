package collectors

import (
	"context"

	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
	"github.com/hydn-co/mesh-slack/internal/helpers"
	"github.com/hydn-co/mesh-slack/internal/options"
	slackapi "github.com/hydn-co/mesh-slack/internal/slack_api"
)

// SlackUsersCollector collects users from Slack workspaces and emits them as
// catalog entities.
type SlackUsersCollector struct {
	*connector.TypedFeatureContext[*options.SlackUsersCollectorOptions]
	initialized bool
}

// NewSlackUsersCollector constructs a SlackUsersCollector.
func NewSlackUsersCollector(ctx *connector.TypedFeatureContext[*options.SlackUsersCollectorOptions]) runner.Feature {
	return &SlackUsersCollector{TypedFeatureContext: ctx}
}

// Init prepares the collector for operation.
func (c *SlackUsersCollector) Init(ctx context.Context) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	c.initialized = true
	return nil
}

// Start begins collecting users from the Slack workspace.
func (c *SlackUsersCollector) Start(ctx context.Context) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	if err := helpers.CheckInitialized(c.initialized); err != nil {
		return err
	}

	return nil
}

// Stop halts user collection and releases resources.
func (c *SlackUsersCollector) Stop(ctx context.Context) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	if err := helpers.CheckInitialized(c.initialized); err != nil {
		return err
	}

	c.initialized = false
	return nil
}
