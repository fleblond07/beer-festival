package main

import (
	"testing"
)

func TestNewDatabase(t *testing.T) {
	t.Run("returns error when URL is empty", func(t *testing.T) {
		_, err := NewDatabase("", "test-key")

		if err == nil {
			t.Error("Expected error when URL is empty, got nil")
		}

		if err.Error() != "SUPABASE_URL and SUPABASE_KEY environment variables are required" {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	t.Run("returns error when key is empty", func(t *testing.T) {
		_, err := NewDatabase("https://test.supabase.co", "")

		if err == nil {
			t.Error("Expected error when key is empty, got nil")
		}

		if err.Error() != "SUPABASE_URL and SUPABASE_KEY environment variables are required" {
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
