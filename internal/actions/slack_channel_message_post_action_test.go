package actions

import (
	"context"
	"testing"

	"github.com/fgrzl/json/polymorphic"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-slack/internal/options"
	"github.com/hydn-co/mesh-slack/internal/payloads"
)

func newTestAction(opts *options.SlackChannelMessagePostActionOptions, p *payloads.SlackChannelMessagePostPayload) *SlackChannelMessagePostAction {
	cfg := &connector.Configuration{
		Options: polymorphic.NewEnvelope(opts),
	}
	if p != nil {
		cfg.Payload = polymorphic.NewEnvelope(p)
	}
	ctx := connector.NewTypedFeatureContext[
		*options.SlackChannelMessagePostActionOptions,
		*payloads.SlackChannelMessagePostPayload,
	](connector.NewFeatureContext(connector.WithConfiguration(cfg)))
	return &SlackChannelMessagePostAction{TypedFeatureContext: ctx}
}

func TestShouldRejectInitWhenChannelMessagePayloadMissing(t *testing.T) {
	// Arrange
	action := newTestAction(&options.SlackChannelMessagePostActionOptions{ChannelID: "C123"}, nil)

	// Act
	err := action.Init(context.Background())

	// Assert
	if err == nil || err.Error() != "message payload is required" {
		t.Fatalf("expected missing payload error, got %v", err)
	}
}

func TestShouldRejectInitWhenChannelMessageIsEmpty(t *testing.T) {
	// Arrange
	action := newTestAction(
		&options.SlackChannelMessagePostActionOptions{ChannelID: "C123"},
		&payloads.SlackChannelMessagePostPayload{Message: "   "},
	)

	// Act
	err := action.Init(context.Background())

	// Assert
	if err == nil || err.Error() != "message cannot be empty" {
		t.Fatalf("expected empty message error, got %v", err)
	}
}

func TestShouldTrimMessageWhenWhitespacePresent(t *testing.T) {
	// Arrange
	input := "  hello slack  "

	// Act
	got, err := verifyMessage(input)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello slack" {
		t.Fatalf("expected trimmed message, got %q", got)
	}
}

func TestShouldRejectMessageWhenEmpty(t *testing.T) {
	// Act
	_, err := verifyMessage("")

	// Assert
	if err == nil || err.Error() != "message cannot be empty" {
		t.Fatalf("expected empty error, got %v", err)
	}
}

func TestShouldRejectMessageWhenTooLong(t *testing.T) {
	// Arrange
	long := make([]byte, 4001)
	for i := range long {
		long[i] = 'a'
	}

	// Act
	_, err := verifyMessage(string(long))

	// Assert
	if err == nil {
		t.Fatal("expected too-long error")
	}
}
