package credentials

import (
	"encoding/json"
	"fmt"
	"strings"
)

type apiKeyCredentials struct {
	APIKey string `json:"api_key"`
}

// ExtractToken returns the bot token from the registered api_key credential field.
func ExtractToken(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "", fmt.Errorf("api key credentials are required")
	}

	var creds apiKeyCredentials
	if err := json.Unmarshal(raw, &creds); err != nil {
		return "", fmt.Errorf("decode api key credentials: %w", err)
	}

	token := strings.TrimSpace(creds.APIKey)
	if token == "" {
		return "", fmt.Errorf("api_key is required")
	}

	return token, nil
}
