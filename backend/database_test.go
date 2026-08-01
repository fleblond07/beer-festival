package main

import (
	"strings"
	"testing"
	"time"
)

func TestNewDatabase(t *testing.T) {
	t.Run("returns error when URL is empty", func(t *testing.T) {
		_, err := NewDatabase("", "test-secret")

		if err == nil {
			t.Error("Expected error when URL is empty, got nil")
		}

		if err.Error() != "DATABASE_URL and JWT_SECRET environment variables are required" {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	t.Run("returns error when key is empty", func(t *testing.T) {
		_, err := NewDatabase("postgres://test:test@localhost:5432/test?sslmode=disable", "")

		if err == nil {
			t.Error("Expected error when key is empty, got nil")
		}

		if err.Error() != "DATABASE_URL and JWT_SECRET environment variables are required" {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	t.Run("returns error when both URL and key are empty", func(t *testing.T) {
		_, err := NewDatabase("", "")

		if err == nil {
			t.Error("Expected error when both are empty, got nil")
		}
	})
}

func TestDatabaseTokenVerification(t *testing.T) {
	db := &Database{jwtSecret: []byte("test-secret")}
	user := User{ID: "user-1", Email: "admin@example.com"}

	t.Run("verifies a signed token", func(t *testing.T) {
		token, err := db.signToken(user, time.Hour)
		if err != nil {
			t.Fatalf("Expected token to be signed, got error: %v", err)
		}

		verifiedUser, err := db.VerifyToken(token)
		if err != nil {
			t.Fatalf("Expected token to verify, got error: %v", err)
		}

		if verifiedUser.ID != user.ID || verifiedUser.Email != user.Email {
			t.Errorf("Expected verified user %+v, got %+v", user, verifiedUser)
		}
	})

	t.Run("rejects a tampered token", func(t *testing.T) {
		token, err := db.signToken(user, time.Hour)
		if err != nil {
			t.Fatalf("Expected token to be signed, got error: %v", err)
		}

		parts := strings.Split(token, ".")
		parts[1] = "tampered"
		_, err = db.VerifyToken(strings.Join(parts, "."))

		if err == nil {
			t.Error("Expected tampered token to be rejected")
		}
	})

	t.Run("rejects an expired token", func(t *testing.T) {
		token, err := db.signToken(user, -time.Hour)
		if err != nil {
			t.Fatalf("Expected token to be signed, got error: %v", err)
		}

		_, err = db.VerifyToken(token)
		if err == nil {
			t.Error("Expected expired token to be rejected")
		}
	})
}
