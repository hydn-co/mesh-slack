package options

import (
	"encoding/json"
	"fmt"

	"github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"
)

// SlackUserMessagePostActionOptions configures the Slack user DM message post action.
type SlackUserMessagePostActionOptions struct {
	// Emails contains one to eight recipient email addresses. A single email
	// opens a 1:1 DM; two to eight open a group DM (MPIM).
	Emails []string `json:"emails" title:"Recipient Email(s)" description:"One to eight recipient email addresses used to open a direct or group message" binding:"required" x-lookup:"{\"entity-type\": \"accounts\", \"display-key\": \"primary_email.address\", \"submit-key\": \"primary_email.address\", \"form-input-type\": \"multi-select\"}"`
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

	var emails []string
	if err := json.Unmarshal(raw.Emails, &emails); err == nil {
		o.Emails = emails
		return nil
	}

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
