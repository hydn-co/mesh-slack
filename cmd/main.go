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
		"Collect Users",
		"Collects users from Slack workspaces.",
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
		"Collect Channels",
		"Collects channels from Slack workspaces.",
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
		"Send Message To Channel",
		"Posts a message to a Slack channel.",
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
		"Send Message To Users",
		"Posts a direct message to a single user chat or up to 8 users in a group chat.",
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
