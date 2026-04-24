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
