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

// SlackUsersCollectorOptions configures the Slack users collector.
type SlackUsersCollectorOptions struct {
	// IncludeWorkspaceUser optionally filters to include or exclude the workspace user.
	IncludeWorkspaceUser *bool `json:"include_workspace_user"`
	// IncludeDisabledUsers optionally filters to include or exclude disabled users.
	IncludeDisabledUsers *bool `json:"include_disabled_users"`
	// MaxResults optionally limits the number of users collected.
	MaxResults *int64 `json:"max_results"`
}

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
type SlackChannelsCollectorOptions struct {
	// IncludeArchived optionally filters to include or exclude archived channels.
	IncludeArchived *bool `json:"include_archived"`
	// IncludePrivate optionally filters to include or exclude private channels.
	IncludePrivate *bool `json:"include_private"`
	// MaxResults optionally limits the number of channels collected.
	MaxResults *int64 `json:"max_results"`
}

func (o *SlackChannelsCollectorOptions) GetDiscriminator() string {
	return "mesh://slack/channels_collector_options"
}

func (o *SlackChannelsCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Activity}
}

func (o *SlackChannelsCollectorOptions) GetRequirements() []string {
	return []string{"slack"}
}
