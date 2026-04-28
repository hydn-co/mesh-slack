package options

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
