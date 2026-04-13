package main

import (
	"github.com/hydn-co/mesh-sdk/pkg/runner"
	"github.com/hydn-co/mesh-slack/internal/slack"
)

func main() {
	manifest := runner.CreateManifest(
		"slack",
		"",
		"Slack",
		"A mesh collector for Slack integration.",
	)

	manifest.RegisterFeature(
		"slack_channel_message_collector",
		"Slack Channel Message Collector",
		"Collects messages from Slack channels and emits them as catalog entities.",
		runner.Schema[*slack.SlackChannelMessageCollectorOptions](),
		runner.GetRequirements[*slack.SlackChannelMessageCollectorOptions](),
		runner.APIKeyCredential,
		runner.Factory[*slack.SlackChannelMessageCollectorOptions](slack.NewSlackChannelMessageCollector),
	)

	manifest.RegisterFeature(
		"slack_channel_message_post_provisioner",
		"Slack Channel Message Post Provisioner",
		"Posts messages to Slack channels based on catalog events.",
		runner.Schema[*slack.SlackChannelMessagePostProvisionerOptions](),
		runner.GetRequirements[*slack.SlackChannelMessagePostProvisionerOptions](),
		runner.APIKeyCredential,
		runner.Factory[*slack.SlackChannelMessagePostProvisionerOptions](slack.NewSlackChannelMessagePostProvisioner),
	)

	runner.Run(manifest)
}
