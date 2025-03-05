package auth

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestMakeJWTAndValidateJWT(t *testing.T) {
	// Shared constants for testing
	tokenSecret := "supersecretkey"
	userID := uuid.New() // Generate a new UUID for testing

	// Test JWT creation
	token, err := MakeJWT(userID, tokenSecret)
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

func TestValidateJWTWithWrongSecret(t *testing.T) {
	originalSecret := "supersecretkey"
	wrongSecret := "wrongsupersecret"
	userID := uuid.New()

	// Create a token with the correct secret
	token, err := MakeJWT(userID, originalSecret)
	if err != nil {
		t.Fatalf("failed to create JWT: %s", err)
	}

	// Attempt to validate the token with the wrong secret
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("expected error for token signed with wrong secret but got none")
	}
}

func TestGetBearerToken(t *testing.T) {
	// Test valid bearer token
	t.Run("Valid Bearer Token", func(t *testing.T) {
		headers := http.Header{}
		headers.Add("Authorization", "Bearer abc123.def456.ghi789")

		token, err := GetBearerToken(headers)
		if err != nil {
			t.Fatalf("expected no error but got: %v", err)
		}

		if token != "abc123.def456.ghi789" {
			t.Errorf("expected token 'abc123.def456.ghi789', got '%s'", token)
		}
	})

	// Test missing Authorization header
	t.Run("Missing Authorization Header", func(t *testing.T) {
		headers := http.Header{}

		_, err := GetBearerToken(headers)
		if err == nil {
			t.Error("expected error for missing header but got none")
		}
	})

	// Test empty Authorization header
	t.Run("Empty Authorization Header", func(t *testing.T) {
		headers := http.Header{}
		headers.Add("Authorization", "")

		_, err := GetBearerToken(headers)
		if err == nil {
			t.Error("expected error for empty header but got none")
		}
	})

	// Test Authorization header without Bearer prefix
	t.Run("No Bearer Prefix", func(t *testing.T) {
		headers := http.Header{}
		headers.Add("Authorization", "abc123.def456.ghi789")

		_, err := GetBearerToken(headers)
		if err == nil {
			t.Error("expected error for missing Bearer prefix but got none")
		}
	})

	// Test Authorization header with extra spaces
	t.Run("Extra Spaces in Header", func(t *testing.T) {
		headers := http.Header{}
		headers.Add("Authorization", "Bearer     abc123.def456.ghi789")

		token, err := GetBearerToken(headers)
		if err != nil {
			t.Fatalf("expected no error but got: %v", err)
		}

		if token != "abc123.def456.ghi789" {
			t.Errorf("expected token 'abc123.def456.ghi789', got '%s'", token)
		}
	})

	// Test case-insensitivity of the Bearer prefix
	t.Run("Case Insensitive Bearer", func(t *testing.T) {
		headers := http.Header{}
		headers.Add("Authorization", "bearer abc123.def456.ghi789")

		token, err := GetBearerToken(headers)
		if err != nil {
			t.Fatalf("expected no error but got: %v", err)
		}

		if token != "abc123.def456.ghi789" {
			t.Errorf("expected token 'abc123.def456.ghi789', got '%s'", token)
		}
	})
}
