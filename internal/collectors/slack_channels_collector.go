package collectors

import (
	"context"

	"github.com/fgrzl/json/polymorphic"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func init() {
	polymorphic.RegisterType[SlackChannelsCollectorOptions]()
}

// SlackChannelsCollectorOptions configures the Slack channels collector.
type SlackChannelsCollectorOptions struct {
	// IncludeArchived optionally filters to include or exclude archived channels.
	IncludeArchived *bool `json:"include_archived"`
	// IncludePrivate optionally filters to include or exclude private channels.
	IncludePrivate *bool `json:"include_private"`
	// MaxResults optionally limits the number of channels collected.
	MaxResults *int64 `json:"max_results"`
}

func (o *SlackChannelsCollectorOptions) GetDiscriminator() string {
	return "mesh://slack/channels_collector_options"
}

func (o *SlackChannelsCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Activity}
}

func (o *SlackChannelsCollectorOptions) GetRequirements() []string {
	return []string{"slack"}
}

// SlackChannelsCollector collects channels from Slack workspaces and emits them
// as catalog entities.
type SlackChannelsCollector struct {
	ctx *connector.TypedFeatureContext[*SlackChannelsCollectorOptions]
}

// NewSlackChannelsCollector constructs a SlackChannelsCollector.
func NewSlackChannelsCollector(ctx *connector.TypedFeatureContext[*SlackChannelsCollectorOptions]) runner.Feature {
	return &SlackChannelsCollector{ctx: ctx}
}

// Init prepares the collector for operation.
func (c *SlackChannelsCollector) Init(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return nil
}

// Start begins collecting channels from the Slack workspace.
func (c *SlackChannelsCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return nil
}

// Stop halts channel collection and releases resources.
func (c *SlackChannelsCollector) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return nil
}
