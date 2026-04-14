package main

import (
	"github.com/hydn-co/mesh-sdk/pkg/runner"
	"github.com/hydn-co/mesh-slack/internal/actions"
	"github.com/hydn-co/mesh-slack/internal/collectors"
	"github.com/hydn-co/mesh-slack/internal/options"
)

func main() {
	runner.Run(WithManifest())
}

func WithManifest() *runner.Manifest {
	manifest := runner.CreateManifest(
		"mesh-slack",
		"",
		"Slack",
		"Mesh integration with Slack",
	)

	// Register Collectors
	manifest.RegisterFeature(
		"slack_users_collector",
		"Slack Users Collector",
		"Collects users from Slack workspaces and emits them as catalog entities.",
		true,
		runner.FeatureTypeCollector,
		runner.Schema[*options.SlackUsersCollectorOptions](),
		runner.GetRequirements[*options.SlackUsersCollectorOptions](),
		runner.APIKeyCredential,
		runner.Factory(collectors.NewSlackUsersCollector),
	)

	manifest.RegisterFeature(
		"slack_channels_collector",
		"Slack Channels Collector",
		"Collects channels from Slack workspaces and emits them as catalog entities.",
		true,
		runner.FeatureTypeCollector,
		runner.Schema[*options.SlackChannelsCollectorOptions](),
		runner.GetRequirements[*options.SlackChannelsCollectorOptions](),
		runner.APIKeyCredential,
		runner.Factory(collectors.NewSlackChannelsCollector),
	)

	// Register Actions
	manifest.RegisterFeature(
		"slack_channel_message_post_action",
		"Slack Channel Message Post Action",
		"Posts messages to Slack channels based on catalog events.",
		false,
		runner.FeatureTypeAction,
		runner.Schema[*options.SlackChannelMessagePostActionOptions](),
		runner.GetRequirements[*options.SlackChannelMessagePostActionOptions](),
		runner.APIKeyCredential,
		runner.Factory(actions.NewSlackChannelMessagePostAction),
	)

	return manifest
}
