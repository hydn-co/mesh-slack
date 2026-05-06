package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

// SlackChannelsCollectorOptions configures the Slack channels collector.
type SlackChannelsCollectorOptions struct{}

func (o *SlackChannelsCollectorOptions) GetDiscriminator() string {
	return "mesh://slack/channels_collector_options"
}

func (o *SlackChannelsCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Activity}
}

func (o *SlackChannelsCollectorOptions) GetRequirements() []string {
	return []string{"slack"}
}
