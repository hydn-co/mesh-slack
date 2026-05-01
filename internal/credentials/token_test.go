package credentials

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldExtractTokenWhenAPIKeyProvided(t *testing.T) {
	token, err := ExtractToken(json.RawMessage(`{"api_key":"xoxb-123"}`))

	require.NoError(t, err)
	assert.Equal(t, "xoxb-123", token)
}

func TestShouldTrimWhitespaceWhenExtractingToken(t *testing.T) {
	token, err := ExtractToken(json.RawMessage(`{"api_key":" xoxb-123 "}`))

	require.NoError(t, err)
	assert.Equal(t, "xoxb-123", token)
}

func TestShouldRejectCredentialsWhenAPIKeyMissing(t *testing.T) {
	_, err := ExtractToken(json.RawMessage(`{"token":"xoxb-123"}`))

	require.Error(t, err)
	assert.EqualError(t, err, "api_key is required")
}

func TestShouldRejectCredentialsWhenJSONInvalid(t *testing.T) {
	_, err := ExtractToken(json.RawMessage(`{"api_key":`))

	require.Error(t, err)
	assert.ErrorContains(t, err, "decode api key credentials")
}
