package auth

import (
	"testing"
)

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	password := "my_secure_password"

	// Test HashPassword
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an unexpected error: %s", err)
	}

	// Ensure the hash is not empty
	if hash == "" {
		t.Fatal("Expected a non-empty hash from HashPassword")
	}

	// Test CheckPasswordHash with matching passwords
	err = CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash failed with matching password: %s", err)
	}

	// Test CheckPasswordHash with non-matching passwords
	err = CheckPasswordHash("wrong_password", hash)
	if err == nil {
		t.Fatal("CheckPasswordHash should have failed with a non-matching password but did not")
	}
}
