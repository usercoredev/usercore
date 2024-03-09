package utils

import (
	"testing"
)

func TestGeneratePasswordHashAndCheckPasswordHash(t *testing.T) {
	password := "secretPassword123"

	// Test GeneratePasswordHash
	hash, err := GeneratePasswordHash(password)
	if err != nil {
		t.Errorf("Failed to generate password hash: %v", err)
	}

	if hash == "" {
		t.Error("Generated hash is empty")
	}

	// Test CheckPasswordHash with correct password
	match := CheckPasswordHash(password, hash)
	if !match {
		t.Errorf("Password and hash should match, but they don't")
	}

	// Test CheckPasswordHash with incorrect password
	wrongPassword := "incorrectPassword"
	match = CheckPasswordHash(wrongPassword, hash)
	if match {
		t.Errorf("Password and hash should not match, but they do")
	}
}
