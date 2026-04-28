package payloads

import "github.com/fgrzl/json/polymorphic"

func init() {
	polymorphic.RegisterType[SlackChannelMessagePostPayload]()
}

// SlackChannelMessagePostPayload is the action payload schema for posting a message.
type SlackChannelMessagePostPayload struct {
	Message string `json:"message" binding:"required"`
}

func (p *SlackChannelMessagePostPayload) GetDiscriminator() string {
	return "mesh://payloads/slack-channel-message-post"
}
