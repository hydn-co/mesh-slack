package channels

import (
	"context"

	"github.com/fgrzl/enumerators"
)

type pageFetcher[T any] func(ctx context.Context, cursor string) ([]T, string, error)

type cursorEnumerator[T any] struct {
	ctx     context.Context
	fetch   pageFetcher[T]
	cursor  string
	items   []T
	index   int
	done    bool
	err     error
	current T
}

// ChannelEnumerator paginates through all channels accessible by the bot token.
func ChannelEnumerator(ctx context.Context, token string) enumerators.Enumerator[SlackChannel] {
	return &cursorEnumerator[SlackChannel]{
		ctx: ctx,
		fetch: func(ctx context.Context, cursor string) ([]SlackChannel, string, error) {
			result, err := ListChannels(ctx, token, cursor)
			if err != nil {
				return nil, "", err
			}

			return result.Channels, result.ResponseMetadata.NextCursor, nil
		},
	}
}

// MemberEnumerator paginates through all members of a channel.
func MemberEnumerator(ctx context.Context, token, channelID string) enumerators.Enumerator[string] {
	return &cursorEnumerator[string]{
		ctx: ctx,
		fetch: func(ctx context.Context, cursor string) ([]string, string, error) {
			result, err := ListMembers(ctx, token, channelID, cursor)
			if err != nil {
				return nil, "", err
			}

			return result.Members, result.ResponseMetadata.NextCursor, nil
		},
	}
}

func (e *cursorEnumerator[T]) MoveNext() bool {
	if e.err != nil || e.done {
		return false
	}

	for e.index >= len(e.items) {
		if err := e.ctx.Err(); err != nil {
			e.err = err
			return false
		}

		items, nextCursor, err := e.fetch(e.ctx, e.cursor)
		if err != nil {
			e.err = err
			return false
		}

		e.items = items
		e.index = 0
		e.cursor = nextCursor
		if len(e.items) == 0 {
			if e.cursor == "" {
				e.done = true
				return false
			}
			continue
		}
	}

	e.current = e.items[e.index]
	e.index++

	if e.index >= len(e.items) && e.cursor == "" {
		e.done = true
	}

	return true
}

func (e *cursorEnumerator[T]) Current() (T, error) {
	if e.err != nil {
		var zero T
		return zero, e.err
	}

	return e.current, nil
}

func (e *cursorEnumerator[T]) Err() error {
	return e.err
}

func (e *cursorEnumerator[T]) Dispose() {
	e.items = nil
}
