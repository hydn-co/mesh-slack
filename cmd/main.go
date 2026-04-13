package main

import (
	"github.com/hydn-co/mesh-sdk/pkg/runner"
	"github.com/hydn-co/mesh-slack/internal/collectors"
	"github.com/hydn-co/mesh-slack/internal/provisioners"
)

func main() {
	runner.Run(WithManifest())
}

func WithManifest() *runner.Manifest {
	manifest := runner.CreateManifest(
		"slack",
		"",
		"Slack",
		"A mesh integration for Slack.",
	)

	manifest.RegisterFeature(
		"slack_channel_message_collector",
		"Slack Channel Message Collector",
		"Collects messages from Slack channels and emits them as catalog entities.",
		runner.Schema[*collectors.SlackChannelMessageCollectorOptions](),
		runner.GetRequirements[*collectors.SlackChannelMessageCollectorOptions](),
		runner.APIKeyCredential,
		runner.Factory[*collectors.SlackChannelMessageCollectorOptions](collectors.NewSlackChannelMessageCollector),
	)

	manifest.RegisterFeature(
		"slack_channel_message_post_provisioner",
		"Slack Channel Message Post Provisioner",
		"Posts messages to Slack channels based on catalog events.",
		runner.Schema[*provisioners.SlackChannelMessagePostProvisionerOptions](),
		runner.GetRequirements[*provisioners.SlackChannelMessagePostProvisionerOptions](),
		runner.APIKeyCredential,
		runner.Factory[*provisioners.SlackChannelMessagePostProvisionerOptions](provisioners.NewSlackChannelMessagePostProvisioner),
	)

	return manifest
}
