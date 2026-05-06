package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

// SlackUsersCollectorOptions configures the Slack users collector.
type SlackUsersCollectorOptions struct{}

func (o *SlackUsersCollectorOptions) GetDiscriminator() string {
	return "mesh://slack/users_collector_options"
}

func (o *SlackUsersCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Activity}
}

func (o *SlackUsersCollectorOptions) GetRequirements() []string {
	return []string{"slack"}
}
