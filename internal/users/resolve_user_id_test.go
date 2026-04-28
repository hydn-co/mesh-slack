package users

import (
	"context"
	"testing"
)

func TestShouldRejectEmptyEmailWhenFindingUserID(t *testing.T) {
	// Arrange
	email := "  "

	// Act
	_, err := FindUserIDByEmail(context.Background(), "token", email)

	// Assert
	if err == nil || err.Error() != "email is required" {
		t.Fatalf("expected email validation error, got %v", err)
	}
}

func TestShouldRejectEmptyEmailEntryWhenResolvingMultiple(t *testing.T) {
	// Arrange
	emails := []string{"a@example.com", "  "}

	// Act
	_, err := ResolveUserIDsByEmails(context.Background(), "token", emails)

	// Assert
	if err == nil || err.Error() != "email list contains an empty entry" {
		t.Fatalf("expected empty entry error, got %v", err)
	}
}

func TestShouldReturnEmptyImmediatelyWhenNoEmailsProvided(t *testing.T) {
	// Arrange: no emails — function must return without any API call.
	// (Token is intentionally invalid; a real API call would fail.)

	// Act
	ids, err := ResolveUserIDsByEmails(context.Background(), "invalid-token", []string{})

	// Assert
	if err != nil {
		t.Fatalf("expected no error for empty input, got %v", err)
	}
	if len(ids) != 0 {
		t.Fatalf("expected empty result, got %v", ids)
	}
}
