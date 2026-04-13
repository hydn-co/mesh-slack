package provisioners

import (
	"context"

	"github.com/fgrzl/json/polymorphic"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func init() {
	polymorphic.RegisterType[SlackChannelMessagePostProvisionerOptions]()
}

// SlackChannelMessagePostProvisionerOptions configures the Slack channel message
// post provisioner.
type SlackChannelMessagePostProvisionerOptions struct {
	// ChannelID is the Slack channel ID to post messages to.
	ChannelID string `json:"channel_id"`
}

func (o *SlackChannelMessagePostProvisionerOptions) GetDiscriminator() string {
	return "mesh://slack/channel_message_post_provisioner_options"
}

func (o *SlackChannelMessagePostProvisionerOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Activity}
}

func (o *SlackChannelMessagePostProvisionerOptions) GetRequirements() []string {
	return []string{"Slack"}
}

// SlackChannelMessagePostProvisioner posts messages to Slack channels based on
// catalog events.
type SlackChannelMessagePostProvisioner struct {
	ctx *connector.TypedFeatureContext[*SlackChannelMessagePostProvisionerOptions]
}

// NewSlackChannelMessagePostProvisioner constructs a SlackChannelMessagePostProvisioner.
func NewSlackChannelMessagePostProvisioner(ctx *connector.TypedFeatureContext[*SlackChannelMessagePostProvisionerOptions]) runner.Feature {
	return &SlackChannelMessagePostProvisioner{ctx: ctx}
}

// Init prepares the provisioner for operation.
func (p *SlackChannelMessagePostProvisioner) Init(_ context.Context) error {
	return nil
}

// Start begins processing catalog events and posting messages to Slack.
func (p *SlackChannelMessagePostProvisioner) Start(_ context.Context) error {
	return nil
}

// Stop halts event processing and releases resources.
func (p *SlackChannelMessagePostProvisioner) Stop(_ context.Context) error {
	return nil
}
