package payloads

// SlackChannelMessagePostPayload is the action payload schema for posting a message.
type SlackChannelMessagePostPayload struct {
	Message string `json:"message"`
}

func (p *SlackChannelMessagePostPayload) GetDiscriminator() string {
	return "mesh://payloads/slack-channel-message-post"
}
