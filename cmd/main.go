package main

import (
	"log"

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
	if err := manifest.RegisterFeature(
		"slack_users_collector",
		"Slack Users Collector",
		"Collects users from Slack workspaces and emits them as catalog entities.",
		true,
		runner.FeatureTypeCollector,
		new(options.SlackUsersCollectorOptions),
		nil,
		runner.APIKeyCredential,
		runner.Factory(collectors.NewSlackUsersCollector),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.RegisterFeature(
		"slack_channels_collector",
		"Slack Channels Collector",
		"Collects channels from Slack workspaces and emits them as catalog entities.",
		true,
		runner.FeatureTypeCollector,
		new(options.SlackChannelsCollectorOptions),
		nil,
		runner.APIKeyCredential,
		runner.Factory(collectors.NewSlackChannelsCollector),
	); err != nil {
		log.Fatal(err)
	}

	// Register Actions
	if err := manifest.RegisterFeature(
		"slack_channel_message_post_action",
		"Slack Channel Message Post Action",
		"Posts messages to Slack channels based on catalog events.",
		false,
		runner.FeatureTypeAction,
		new(options.SlackChannelMessagePostActionOptions),
		new(payloads.SlackChannelMessagePostPayload),
		runner.APIKeyCredential,
		runner.Factory(actions.NewSlackChannelMessagePostAction),
	); err != nil {
		log.Fatal(err)
	}

	err := manifest.Validate(); if err != nil {
		log.Fatal(err)
	}

	return manifest
}
