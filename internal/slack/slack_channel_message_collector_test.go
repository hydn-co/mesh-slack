package slack_test

import (
	"context"
	"testing"

	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-slack/internal/slack"
	"github.com/stretchr/testify/assert"
)

func TestSlackChannelMessageCollectorShouldReturnCorrectDiscriminator(t *testing.T) {
	opts := &slack.SlackChannelMessageCollectorOptions{}
	assert.Equal(t, "mesh://slack/channel_message_collector_options", opts.GetDiscriminator())
}

func TestSlackChannelMessageCollectorShouldReturnActivitySpace(t *testing.T) {
	opts := &slack.SlackChannelMessageCollectorOptions{}
	spaces := opts.GetSpaces()
	assert.Contains(t, spaces, "activity")
}

func TestSlackChannelMessageCollectorShouldReturnSlackRequirement(t *testing.T) {
	opts := &slack.SlackChannelMessageCollectorOptions{}
	assert.Contains(t, opts.GetRequirements(), "Slack")
}

func TestSlackChannelMessageCollectorInitShouldSucceed(t *testing.T) {
	ctx := connector.NewTypedFeatureContext[*slack.SlackChannelMessageCollectorOptions](
		connector.NewFeatureContext(),
	)
	feature := slack.NewSlackChannelMessageCollector(ctx)

	err := feature.Init(context.Background())
	assert.NoError(t, err)
}

func TestSlackChannelMessageCollectorStartShouldSucceed(t *testing.T) {
	ctx := connector.NewTypedFeatureContext[*slack.SlackChannelMessageCollectorOptions](
		connector.NewFeatureContext(),
	)
	feature := slack.NewSlackChannelMessageCollector(ctx)

	err := feature.Start(context.Background())
	assert.NoError(t, err)
}

func TestSlackChannelMessageCollectorStopShouldSucceed(t *testing.T) {
	ctx := connector.NewTypedFeatureContext[*slack.SlackChannelMessageCollectorOptions](
		connector.NewFeatureContext(),
	)
	feature := slack.NewSlackChannelMessageCollector(ctx)

	err := feature.Stop(context.Background())
	assert.NoError(t, err)
}
