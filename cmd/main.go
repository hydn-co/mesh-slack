package main

import (
	"log"

	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
	"github.com/hydn-co/mesh-slack/internal/actions"
	"github.com/hydn-co/mesh-slack/internal/collectors"
	"github.com/hydn-co/mesh-slack/internal/options"
	"github.com/hydn-co/mesh-slack/internal/payloads"
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
	manifest.MustRegisterFeature(
		"slack_users_collector",
		"Slack Users Collector",
		"Collects users from Slack workspaces and emits them as catalog entities.",
		runner.FeatureSchedulable,
		runner.FeatureTypeCollector,
		new(options.SlackUsersCollectorOptions),
		(*connector.NoPayload)(nil),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(collectors.NewSlackUsersCollector),
	)

	manifest.MustRegisterFeature(
		"slack_channels_collector",
		"Slack Channels Collector",
		"Collects channels from Slack workspaces and emits them as catalog entities.",
		runner.FeatureSchedulable,
		runner.FeatureTypeCollector,
		new(options.SlackChannelsCollectorOptions),
		(*connector.NoPayload)(nil),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(collectors.NewSlackChannelsCollector),
	)

	// Register Actions
	manifest.MustRegisterFeature(
		"slack_channel_message_post_action",
		"Slack Channel Message Post Action",
		"Posts messages to Slack channels based on catalog events.",
		runner.FeatureUnschedulable,
		runner.FeatureTypeAction,
		new(options.SlackChannelMessagePostActionOptions),
		new(payloads.SlackChannelMessagePostPayload),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(actions.NewSlackChannelMessagePostAction),
	)

	manifest.MustRegisterFeature(
		"slack_user_message_post_action",
		"Slack User Message Post Action",
		"Posts a direct message to a Slack user based on catalog events.",
		runner.FeatureUnschedulable,
		runner.FeatureTypeAction,
		new(options.SlackUserMessagePostActionOptions),
		new(payloads.SlackChannelMessagePostPayload),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(actions.NewSlackUserMessagePostAction),
	)

	if err := manifest.Validate(); err != nil {
		log.Fatal(err)
	}

	return manifest
}
