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

func TestShouldRejectMissingMessagePayload(t *testing.T) {
	action := newTestAction(&options.SlackChannelMessagePostActionOptions{ChannelID: "C123"}, nil)

	err := action.Init(context.Background())
	if err == nil || err.Error() != "message payload is required" {
		t.Fatalf("expected missing payload error, got %v", err)
	}
}

func TestShouldRejectEmptyMessagePayload(t *testing.T) {
	action := newTestAction(
		&options.SlackChannelMessagePostActionOptions{ChannelID: "C123"},
		&payloads.SlackChannelMessagePostPayload{Message: "   "},
	)

	err := action.Init(context.Background())
	if err == nil || err.Error() != "message cannot be empty" {
		t.Fatalf("expected empty message error, got %v", err)
	}
}

func TestShouldVerifyMessageTrimsWhitespace(t *testing.T) {
	got, err := verifyMessage("  hello slack  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello slack" {
		t.Fatalf("expected trimmed message, got %q", got)
	}
}

func TestShouldVerifyMessageRejectsEmpty(t *testing.T) {
	_, err := verifyMessage("")
	if err == nil || err.Error() != "message cannot be empty" {
		t.Fatalf("expected empty error, got %v", err)
	}
}

func TestShouldVerifyMessageRejectsTooLong(t *testing.T) {
	long := make([]byte, 4001)
	for i := range long {
		long[i] = 'a'
	}
	_, err := verifyMessage(string(long))
	if err == nil {
		t.Fatal("expected too-long error")
	}
}
