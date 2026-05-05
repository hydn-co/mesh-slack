package options

import (
	"encoding/json"
	"fmt"

	"github.com/fgrzl/json/polymorphic"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"
)

func init() {
	polymorphic.RegisterType[SlackUsersCollectorOptions]()
	polymorphic.RegisterType[SlackChannelsCollectorOptions]()
	polymorphic.RegisterType[SlackChannelMessagePostActionOptions]()
	polymorphic.RegisterType[SlackUserMessagePostActionOptions]()
}

// SlackChannelMessagePostActionOptions configures the Slack channel message post action.
type SlackChannelMessagePostActionOptions struct {
	// ChannelID is the Slack channel ID to post messages to.
	ChannelID string `json:"channel_id" title:"Channel ID" description:"The Slack channel ID to post the message to" binding:"required"`
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

// SlackUserMessagePostActionOptions configures the Slack user DM message post action.
type SlackUserMessagePostActionOptions struct {
	// Emails contains one to eight recipient email addresses. A single email
	// opens a 1:1 DM; two to eight open a group DM (MPIM).
	Emails []string `json:"emails" title:"Recipient Emails" description:"One to eight recipient email addresses used to open a direct or group message" binding:"required"`
}

// UnmarshalJSON accepts both a JSON string and a JSON array for the emails field,
// allowing options stored as a bare string to be read back correctly.
func (o *SlackUserMessagePostActionOptions) UnmarshalJSON(data []byte) error {
	var raw struct {
		Emails json.RawMessage `json:"emails"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw.Emails) == 0 {
		return nil
	}
	// Try array first.
	var emails []string
	if err := json.Unmarshal(raw.Emails, &emails); err == nil {
		o.Emails = emails
		return nil
	}
	// Fall back to bare string → single-element slice.
	var single string
	if err := json.Unmarshal(raw.Emails, &single); err == nil {
		o.Emails = []string{single}
		return nil
	}
	return fmt.Errorf("emails: expected string or array of strings")
}

func (o *SlackUserMessagePostActionOptions) GetDiscriminator() string {
	return "mesh://slack/user_message_post_action_options"
}

func (o *SlackUserMessagePostActionOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Activity}
}

func (o *SlackUserMessagePostActionOptions) GetRequirements() []string {
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
