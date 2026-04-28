package actions

import (
	"context"
	"testing"

	"github.com/fgrzl/json/polymorphic"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-slack/internal/options"
	"github.com/hydn-co/mesh-slack/internal/payloads"
)

func newTestUserAction(opts *options.SlackUserMessagePostActionOptions, p *payloads.SlackChannelMessagePostPayload) *SlackUserMessagePostAction {
	cfg := &connector.Configuration{
		Options: polymorphic.NewEnvelope(opts),
	}
	if p != nil {
		cfg.Payload = polymorphic.NewEnvelope(p)
	}
	ctx := connector.NewTypedFeatureContext[
		*options.SlackUserMessagePostActionOptions,
		*payloads.SlackChannelMessagePostPayload,
	](connector.NewFeatureContext(connector.WithConfiguration(cfg)))
	return &SlackUserMessagePostAction{TypedFeatureContext: ctx}
}

func TestShouldRejectInitWhenNoEmailsProvided(t *testing.T) {
	// Arrange
	action := newTestUserAction(&options.SlackUserMessagePostActionOptions{Emails: nil}, nil)

	// Act
	err := action.Init(context.Background())

	// Assert
	if err == nil || err.Error() != "at least one email is required" {
		t.Fatalf("expected no-emails error, got %v", err)
	}
}

func TestShouldRejectInitWhenTooManyEmailsProvided(t *testing.T) {
	// Arrange
	emails := make([]string, 9)
	for i := range emails {
		emails[i] = "user@example.com"
	}
	action := newTestUserAction(&options.SlackUserMessagePostActionOptions{Emails: emails}, nil)

	// Act
	err := action.Init(context.Background())

	// Assert
	if err == nil || err.Error() != "at most 8 recipients are supported, got 9" {
		t.Fatalf("expected too-many-recipients error, got %v", err)
	}
}

func TestShouldRejectInitWhenUserMessagePayloadMissing(t *testing.T) {
	// Arrange
	action := newTestUserAction(
		&options.SlackUserMessagePostActionOptions{Emails: []string{"person@example.com"}},
		nil,
	)

	// Act
	err := action.Init(context.Background())

	// Assert
	if err == nil || err.Error() != "message payload is required" {
		t.Fatalf("expected missing payload error, got %v", err)
	}
}

func TestShouldRejectInitWhenUserMessageIsEmpty(t *testing.T) {
	// Arrange
	action := newTestUserAction(
		&options.SlackUserMessagePostActionOptions{Emails: []string{"person@example.com"}},
		&payloads.SlackChannelMessagePostPayload{Message: "  "},
	)

	// Act
	err := action.Init(context.Background())

	// Assert
	if err == nil || err.Error() != "message cannot be empty" {
		t.Fatalf("expected empty message error, got %v", err)
	}
}
