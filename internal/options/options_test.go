package options

import (
	"encoding/json"
	"testing"

	"github.com/fgrzl/json/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldExposeOptionTitlesInGeneratedSchema(t *testing.T) {
	channelSchema := jsonschema.SchemaFrom[SlackChannelMessagePostActionOptions]()

	var channelDecoded map[string]any
	require.NoError(t, json.Unmarshal(channelSchema, &channelDecoded))

	channelProperties, ok := channelDecoded["properties"].(map[string]any)
	require.True(t, ok)

	channelIDSchema, ok := channelProperties["channel_id"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "Channel ID", channelIDSchema["title"])
	require.Equal(t, "The Slack channel ID to post the message to", channelIDSchema["description"])

	userSchema := jsonschema.SchemaFrom[SlackUserMessagePostActionOptions]()

	var userDecoded map[string]any
	require.NoError(t, json.Unmarshal(userSchema, &userDecoded))

	userProperties, ok := userDecoded["properties"].(map[string]any)
	require.True(t, ok)

	emailsSchema, ok := userProperties["emails"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "Recipient Email(s)", emailsSchema["title"])
	require.Equal(t, "One to eight recipient email addresses used to open a direct or group message", emailsSchema["description"])
}

func TestShouldUnmarshalEmailsWhenStoredAsArray(t *testing.T) {
	// Arrange
	data := `{"emails":["a@example.com","b@example.com"]}`

	// Act
	var o SlackUserMessagePostActionOptions
	err := json.Unmarshal([]byte(data), &o)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"a@example.com", "b@example.com"}, o.Emails)
}

func TestShouldUnmarshalEmailsWhenStoredAsBareString(t *testing.T) {
	// Arrange: options persisted before the schema enforced an array
	data := `{"emails":"user@example.com"}`

	// Act
	var o SlackUserMessagePostActionOptions
	err := json.Unmarshal([]byte(data), &o)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"user@example.com"}, o.Emails)
}

func TestShouldMarshalEmailsAsArray(t *testing.T) {
	// Arrange
	o := SlackUserMessagePostActionOptions{Emails: []string{"user@example.com"}}

	// Act
	data, err := json.Marshal(o)

	// Assert
	require.NoError(t, err)
	assert.JSONEq(t, `{"emails":["user@example.com"]}`, string(data))
}
