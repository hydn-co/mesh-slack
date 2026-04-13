package collectors

import (
	"context"

	"github.com/fgrzl/json/polymorphic"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func init() {
	polymorphic.RegisterType[SlackChannelMessageCollectorOptions]()
}

// SlackChannelMessageCollectorOptions configures the Slack channel message collector.
type SlackChannelMessageCollectorOptions struct {
	// ChannelIDs lists the Slack channel IDs to collect messages from.
	ChannelIDs []string `json:"channel_ids"`
}

func (o *SlackChannelMessageCollectorOptions) GetDiscriminator() string {
	return "mesh://slack/channel_message_collector_options"
}

func (o *SlackChannelMessageCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Activity}
}

func (o *SlackChannelMessageCollectorOptions) GetRequirements() []string {
	return []string{"Slack"}
}

// SlackChannelMessageCollector collects messages from Slack channels and emits
// them as catalog entities.
type SlackChannelMessageCollector struct {
	ctx *connector.TypedFeatureContext[*SlackChannelMessageCollectorOptions]
}

// NewSlackChannelMessageCollector constructs a SlackChannelMessageCollector.
func NewSlackChannelMessageCollector(ctx *connector.TypedFeatureContext[*SlackChannelMessageCollectorOptions]) runner.Feature {
	return &SlackChannelMessageCollector{ctx: ctx}
}

// Init prepares the collector for operation.
func (c *SlackChannelMessageCollector) Init(_ context.Context) error {
	return nil
}

// Start begins collecting messages from the configured Slack channels.
func (c *SlackChannelMessageCollector) Start(_ context.Context) error {
	return nil
}

// Stop halts message collection and releases resources.
func (c *SlackChannelMessageCollector) Stop(_ context.Context) error {
	return nil
}
