package channels

import (
	"context"
	"strings"
	"time"

	"github.com/fgrzl/enumerators"
	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
)

type pageFetcher[T any] func(ctx context.Context, cursor string) ([]T, string, error)

const (
	slackThrottleBaseDelay  = 2 * time.Second
	slackThrottleMaxDelay   = 60 * time.Second
	slackThrottleMaxRetries = 5
)

// ChannelEnumerator paginates through all channels accessible by the bot token.
func ChannelEnumerator(ctx context.Context, token string) enumerators.Enumerator[SlackChannel] {
	return cursorEnumerator(ctx, func(ctx context.Context, cursor string) ([]SlackChannel, string, error) {
		result, err := ListChannels(ctx, token, cursor)
		if err != nil {
			return nil, "", err
		}

		return result.Channels, result.ResponseMetadata.NextCursor, nil
	})
}

// MemberEnumerator paginates through all members of a channel.
func MemberEnumerator(ctx context.Context, token, channelID string) enumerators.Enumerator[string] {
	return cursorEnumerator(ctx, func(ctx context.Context, cursor string) ([]string, string, error) {
		result, err := ListMembers(ctx, token, channelID, cursor)
		if err != nil {
			return nil, "", err
		}

		return result.Members, result.ResponseMetadata.NextCursor, nil
	})
}

func cursorEnumerator[T any](
	ctx context.Context,
	fetch pageFetcher[T],
) enumerators.Enumerator[T] {
	cursor := ""

	return connectorutil.ThrottledPageEnumerator(ctx, connectorutil.ThrottlePolicy{
		IsThrottled: isSlackThrottleError,
		BaseDelay:   slackThrottleBaseDelay,
		MaxDelay:    slackThrottleMaxDelay,
		MaxRetries:  slackThrottleMaxRetries,
	}, func() ([]T, bool, error) {
		if err := ctx.Err(); err != nil {
			return nil, false, err
		}

		items, nextCursor, err := fetch(ctx, cursor)
		if err != nil {
			return nil, false, err
		}

		if nextCursor == "" {
			return items, false, nil
		}

		cursor = nextCursor
		return items, true, nil
	})
}

func isSlackThrottleError(err error) bool {
	if err == nil {
		return false
	}

	message := strings.ToLower(err.Error())
	return strings.Contains(message, "rate limit") ||
		strings.Contains(message, "too many requests") ||
		strings.Contains(message, "rate_limited") ||
		strings.Contains(message, "429")
}
