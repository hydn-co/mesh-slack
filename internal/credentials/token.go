package credentials

import (
	"encoding/json"
	"fmt"
)

type tokenCredentials struct {
	Token  string `json:"token"`
	APIKey string `json:"api_key"`
}

// ExtractToken returns the bot token from known credential fields.
func ExtractToken(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "", fmt.Errorf("no credentials provided")
	}

	var creds tokenCredentials
	if err := json.Unmarshal(raw, &creds); err != nil {
		return "", fmt.Errorf("invalid credentials JSON format: %w", err)
	}

	if creds.Token != "" {
		return creds.Token, nil
	}

	if creds.APIKey != "" {
		return creds.APIKey, nil
	}

	return "", fmt.Errorf("missing token credential field (expected token or api_key)")
}
