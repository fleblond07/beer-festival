package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestLoadFestivals(t *testing.T) {
	t.Run("loads festivals from valid JSON file", func(t *testing.T) {
		festivals, err := loadFestivals("./festivals.json")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(festivals) == 0 {
			t.Error("Expected festivals to be loaded, got empty array")
		}

		// Check first festival has required fields
		if festivals[0].ID == "" {
			t.Error("Expected festival to have ID")
		}
		if festivals[0].Name == "" {
			t.Error("Expected festival to have Name")
		}
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		_, err := loadFestivals("./nonexistent.json")
		if err == nil {
			t.Error("Expected error for non-existent file, got nil")
		}
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		// Create a temporary invalid JSON file
		tmpFile, err := os.CreateTemp("", "invalid-*.json")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())

		tmpFile.WriteString("{invalid json}")
		tmpFile.Close()

		_, err = loadFestivals(tmpFile.Name())
		if err == nil {
			t.Error("Expected error for invalid JSON, got nil")
		}
	})
}

func TestGetConfig(t *testing.T) {
	t.Run("uses default values when env vars not set", func(t *testing.T) {
		os.Clearenv()
		config := getConfig()

		if config.Port != "8080" {
			t.Errorf("Expected default port 8080, got %s", config.Port)
		}
		if config.FestivalsFile != "./festivals.json" {
			t.Errorf("Expected default festivals file ./festivals.json, got %s", config.FestivalsFile)
		}
		if config.AllowedOrigins != "*" {
			t.Errorf("Expected default allowed origins *, got %s", config.AllowedOrigins)
		}
	})

	t.Run("uses env vars when set", func(t *testing.T) {
		os.Setenv("PORT", "3000")
		os.Setenv("FESTIVALS_FILE", "./custom.json")
		os.Setenv("ALLOWED_ORIGINS", "http://localhost:5173")
		defer os.Clearenv()

		config := getConfig()

		if config.Port != "3000" {
			t.Errorf("Expected port 3000, got %s", config.Port)
		}
		if config.FestivalsFile != "./custom.json" {
			t.Errorf("Expected festivals file ./custom.json, got %s", config.FestivalsFile)
		}
		if config.AllowedOrigins != "http://localhost:5173" {
			t.Errorf("Expected allowed origins http://localhost:5173, got %s", config.AllowedOrigins)
		}
	})
}

func TestFestivalsHandler(t *testing.T) {
	mockFestivals := []Festival{
		{
			ID:          "1",
			Name:        "Test Festival",
			Description: "A test festival",
			StartDate:   "2025-10-01",
			EndDate:     "2025-10-03",
			City:        "Paris",
			Region:      "Île-de-France",
			Location: Location{
				Latitude:  48.8566,
				Longitude: 2.3522,
			},
			Image:        "https://example.com/image.jpg",
			Website:      "https://example.com",
			BreweryCount: 50,
		},
	}

	t.Run("returns festivals as JSON", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/festivals", nil)
		w := httptest.NewRecorder()

		handler := makeFestivalsHandler(mockFestivals, "*")
		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", contentType)
		}

		var festivals []Festival
		if err := json.NewDecoder(w.Body).Decode(&festivals); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(festivals) != 1 {
			t.Errorf("Expected 1 festival, got %d", len(festivals))
		}

		if festivals[0].ID != "1" {
			t.Errorf("Expected festival ID 1, got %s", festivals[0].ID)
		}
	})

	t.Run("sets CORS headers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/festivals", nil)
		w := httptest.NewRecorder()

		handler := makeFestivalsHandler(mockFestivals, "*")
		handler(w, req)

		origin := w.Header().Get("Access-Control-Allow-Origin")
		if origin != "*" {
			t.Errorf("Expected CORS origin *, got %s", origin)
		}

		methods := w.Header().Get("Access-Control-Allow-Methods")
		if methods != "GET, OPTIONS" {
			t.Errorf("Expected CORS methods 'GET, OPTIONS', got %s", methods)
		}
	})

	t.Run("handles OPTIONS requests", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/api/festivals", nil)
		w := httptest.NewRecorder()

		handler := makeFestivalsHandler(mockFestivals, "*")
		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for OPTIONS, got %d", w.Code)
		}
	})

	t.Run("uses custom allowed origins", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/festivals", nil)
		w := httptest.NewRecorder()

		handler := makeFestivalsHandler(mockFestivals, "http://localhost:5173")
		handler(w, req)

		origin := w.Header().Get("Access-Control-Allow-Origin")
		if origin != "http://localhost:5173" {
			t.Errorf("Expected CORS origin http://localhost:5173, got %s", origin)
		}
	})
}
