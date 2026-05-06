package credentials

import (
	"encoding/json"

	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
)

// ExtractToken returns the bot token from the registered api_key credential field.
func ExtractToken(raw json.RawMessage) (string, error) {
	return connectorutil.ExtractAPIKey(raw)
}
