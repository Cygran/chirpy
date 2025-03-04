package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWTAndValidateJWT(t *testing.T) {
	// Shared constants for testing
	tokenSecret := "supersecretkey"
	expiresIn := 2 * time.Minute
	userID := uuid.New() // Generate a new UUID for testing

	// Test JWT creation
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("failed to create JWT: %s", err)
	}

	// Test validation of the created token
	validatedUUID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("failed to validate JWT: %s", err)
	}

	if validatedUUID != userID {
		t.Errorf("expected userID %v, got %v", userID, validatedUUID)
	}
}

func TestValidateJWTWithInvalidToken(t *testing.T) {
	tokenSecret := "supersecretkey"
	invalidToken := "this.is.not.a.valid.token"

	_, err := ValidateJWT(invalidToken, tokenSecret)
	if err == nil {
		t.Error("expected error for invalid token but got none")
	}
}

func TestValidateJWTWithExpiredToken(t *testing.T) {
	tokenSecret := "supersecretkey"
	expiresIn := -1 * time.Minute // Set token to expire in the past
	userID := uuid.New()

	// Create an expired token
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("failed to create JWT: %s", err)
	}

	// Attempt to validate the expired token
	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Error("expected error for expired token but got none")
	}
}

func TestValidateJWTWithWrongSecret(t *testing.T) {
	originalSecret := "supersecretkey"
	wrongSecret := "wrongsupersecret"
	expiresIn := 2 * time.Minute
	userID := uuid.New()

	// Create a token with the correct secret
	token, err := MakeJWT(userID, originalSecret, expiresIn)
	if err != nil {
		t.Fatalf("failed to create JWT: %s", err)
	}

	// Attempt to validate the token with the wrong secret
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("expected error for token signed with wrong secret but got none")
	}
}
