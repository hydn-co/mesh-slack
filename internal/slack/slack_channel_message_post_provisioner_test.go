package slack_test

import (
	"context"
	"testing"

	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-slack/internal/slack"
	"github.com/stretchr/testify/assert"
)

func TestSlackChannelMessagePostProvisionerShouldReturnCorrectDiscriminator(t *testing.T) {
	opts := &slack.SlackChannelMessagePostProvisionerOptions{}
	assert.Equal(t, "mesh://slack/channel_message_post_provisioner_options", opts.GetDiscriminator())
}

func TestSlackChannelMessagePostProvisionerShouldReturnActivitySpace(t *testing.T) {
	opts := &slack.SlackChannelMessagePostProvisionerOptions{}
	spaces := opts.GetSpaces()
	assert.Contains(t, spaces, "activity")
}

func TestSlackChannelMessagePostProvisionerShouldReturnSlackRequirement(t *testing.T) {
	opts := &slack.SlackChannelMessagePostProvisionerOptions{}
	assert.Contains(t, opts.GetRequirements(), "Slack")
}

func TestSlackChannelMessagePostProvisionerInitShouldSucceed(t *testing.T) {
	ctx := connector.NewTypedFeatureContext[*slack.SlackChannelMessagePostProvisionerOptions](
		connector.NewFeatureContext(),
	)
	feature := slack.NewSlackChannelMessagePostProvisioner(ctx)

	err := feature.Init(context.Background())
	assert.NoError(t, err)
}

func TestSlackChannelMessagePostProvisionerStartShouldSucceed(t *testing.T) {
	ctx := connector.NewTypedFeatureContext[*slack.SlackChannelMessagePostProvisionerOptions](
		connector.NewFeatureContext(),
	)
	feature := slack.NewSlackChannelMessagePostProvisioner(ctx)

	err := feature.Start(context.Background())
	assert.NoError(t, err)
}

func TestSlackChannelMessagePostProvisionerStopShouldSucceed(t *testing.T) {
	ctx := connector.NewTypedFeatureContext[*slack.SlackChannelMessagePostProvisionerOptions](
		connector.NewFeatureContext(),
	)
	feature := slack.NewSlackChannelMessagePostProvisioner(ctx)

	err := feature.Stop(context.Background())
	assert.NoError(t, err)
}
