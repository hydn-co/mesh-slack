package options

import (
	"github.com/fgrzl/json/polymorphic"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"
)

func init() {
	polymorphic.RegisterType[SlackChannelMessagePostActionOptions]()
}

// SlackChannelMessagePostActionOptions configures the Slack channel message post action.
type SlackChannelMessagePostActionOptions struct {
	// ChannelID is the Slack channel ID to post messages to.
	ChannelID string `json:"channel_id" binding:"required"`
	// Message is the required plain-text message posted during this action invocation.
	Message string `json:"message" binding:"required"`
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
