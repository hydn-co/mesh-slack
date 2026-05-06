package payloads

// SlackChannelMessagePostPayload is the action payload schema for posting a message.
type SlackChannelMessagePostPayload struct {
	Message string `json:"message" binding:"required" title:"Channel Message" description:"The message to send to the configured Slack channel"`
}

func (p *SlackChannelMessagePostPayload) GetDiscriminator() string {
	return "mesh://payloads/slack-channel-message-post"
}
