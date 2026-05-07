package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

// SlackChannelMessagePostActionOptions configures the Slack channel message post action.
type SlackChannelMessagePostActionOptions struct {
	// ChannelID is the Slack channel ID to post messages to.
	ChannelID string `json:"channel_id" title:"Channel ID" description:"The Slack channel ID to post the message to" binding:"required"  x-lookup:"{\"entity-type\": \"channels\", \"display-key\": \"name\", \"submit-key\": \"channel_ref\", \"form-input-type\": \"select\"}"`
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
