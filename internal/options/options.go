package options

import (
	"github.com/fgrzl/json/polymorphic"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"
)

func init() {
	polymorphic.RegisterType[SlackUsersCollectorOptions]()
	polymorphic.RegisterType[SlackChannelsCollectorOptions]()
	polymorphic.RegisterType[SlackChannelMessagePostActionOptions]()
}

// SlackChannelMessagePostActionOptions configures the Slack channel message post action.
type SlackChannelMessagePostActionOptions struct {
	// ChannelID is the Slack channel ID to post messages to.
	ChannelID string `json:"channel_id" binding:"required"`
}

func (o *SlackChannelMessagePostActionOptions) GetDiscriminator() string {
	return "mesh://slack/channel_message_post_action_options"
}

func (o *SlackChannelMessagePostActionOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Activity}
}

func (o *SlackChannelMessagePostActionOptions) GetRequirements() []string {
	return []string{"slack"}
}

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
