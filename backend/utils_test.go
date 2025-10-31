package main

import (
	"testing"
	"time"
)

func TestConvertTime(t *testing.T) {
	t.Run("converts valid date string to time.Time", func(t *testing.T) {
		dateStr := "2025-10-01"
		result, err := ConvertTime(dateStr)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		expected := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("returns error for invalid date format", func(t *testing.T) {
		dateStr := "10/01/2025"
		_, err := ConvertTime(dateStr)

		if err == nil {
			t.Error("Expected error for invalid date format, got nil")
		}
	})

	t.Run("returns error for invalid date string", func(t *testing.T) {
		dateStr := "not-a-date"
		_, err := ConvertTime(dateStr)

		if err == nil {
			t.Error("Expected error for invalid date string, got nil")
		}
	})

	t.Run("returns error for empty string", func(t *testing.T) {
		dateStr := ""
		_, err := ConvertTime(dateStr)

		if err == nil {
			t.Error("Expected error for empty string, got nil")
		}
	})

	t.Run("handles leap year dates correctly", func(t *testing.T) {
		dateStr := "2024-02-29"
		result, err := ConvertTime(dateStr)

		if err != nil {
			t.Fatalf("Expected no error for leap year date, got %v", err)
		}

		expected := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("rejects invalid leap year dates", func(t *testing.T) {
		dateStr := "2023-02-29"
		_, err := ConvertTime(dateStr)

		if err == nil {
			t.Error("Expected error for invalid leap year date, got nil")
		}
	})
}
