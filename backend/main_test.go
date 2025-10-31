package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestGetConfig(t *testing.T) {
	t.Run("uses values from environment variables", func(t *testing.T) {
		os.Setenv("PORT", "3000")
		os.Setenv("ALLOWED_ORIGINS", "http://localhost:5173")
		os.Setenv("SUPABASE_URL", "https://test.supabase.co")
		os.Setenv("SUPABASE_KEY", "test-key")
		defer os.Clearenv()

		config := getConfig()

		if config.Port != "3000" {
			t.Errorf("Expected port 3000, got %s", config.Port)
		}
		if config.AllowedOrigins != "http://localhost:5173" {
			t.Errorf("Expected allowed origins http://localhost:5173, got %s", config.AllowedOrigins)
		}
		if config.SupabaseURL != "https://test.supabase.co" {
			t.Errorf("Expected Supabase URL https://test.supabase.co, got %s", config.SupabaseURL)
		}
		if config.SupabaseKey != "test-key" {
			t.Errorf("Expected Supabase key test-key, got %s", config.SupabaseKey)
		}
	})

	t.Run("returns empty strings when env vars not set", func(t *testing.T) {
		os.Clearenv()
		config := getConfig()

		if config.Port != "" {
			t.Errorf("Expected empty port, got %s", config.Port)
		}
		if config.AllowedOrigins != "" {
			t.Errorf("Expected empty allowed origins, got %s", config.AllowedOrigins)
		}
		if config.SupabaseURL != "" {
			t.Errorf("Expected empty Supabase URL, got %s", config.SupabaseURL)
		}
		if config.SupabaseKey != "" {
			t.Errorf("Expected empty Supabase key, got %s", config.SupabaseKey)
		}
	})
}

func TestFestivalsHandler(t *testing.T) {
	mockFestivals := []Festival{
		{
			ID:          1,
			Name:        "Test Festival",
			Description: "A test festival",
			StartDate:   time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2025, 10, 3, 0, 0, 0, 0, time.UTC),
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

		if festivals[0].ID != 1 {
			t.Errorf("Expected festival ID 1, got %d", festivals[0].ID)
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
		req.Header.Set("Origin", "http://localhost:5173")
		w := httptest.NewRecorder()

		handler := makeFestivalsHandler(mockFestivals, "http://localhost:5173")
		handler(w, req)

		origin := w.Header().Get("Access-Control-Allow-Origin")
		if origin != "http://localhost:5173" {
			t.Errorf("Expected CORS origin http://localhost:5173, got %s", origin)
		}
	})

	t.Run("blocks unlisted origins", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/festivals", nil)
		req.Header.Set("Origin", "https://evil.com")
		w := httptest.NewRecorder()

		handler := makeFestivalsHandler(mockFestivals, "http://localhost:5173")
		handler(w, req)

		origin := w.Header().Get("Access-Control-Allow-Origin")
		if origin == "https://evil.com" {
			t.Error("Should not allow unlisted origin")
		}
	})

	t.Run("supports multiple allowed origins", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/festivals", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()

		handler := makeFestivalsHandler(mockFestivals, "http://localhost:5173,http://localhost:3000")
		handler(w, req)

		origin := w.Header().Get("Access-Control-Allow-Origin")
		if origin != "http://localhost:3000" {
			t.Errorf("Expected CORS origin http://localhost:3000, got %s", origin)
		}
	})
}

func TestHealthCheckHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	healthCheckHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status ok, got %s", response["status"])
	}

	if response["version"] == "" {
		t.Error("Expected version to be set")
	}
}

func TestEnableCORS(t *testing.T) {
	t.Run("sets wildcard CORS when allowed origins is *", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/festivals", nil)
		w := httptest.NewRecorder()

		enableCORS(w, req, "*")

		origin := w.Header().Get("Access-Control-Allow-Origin")
		if origin != "*" {
			t.Errorf("Expected CORS origin *, got %s", origin)
		}
	})

	t.Run("sets specific origin when it matches allowed list", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/festivals", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		w := httptest.NewRecorder()

		enableCORS(w, req, "http://localhost:5173")

		origin := w.Header().Get("Access-Control-Allow-Origin")
		if origin != "http://localhost:5173" {
			t.Errorf("Expected CORS origin http://localhost:5173, got %s", origin)
		}
	})

	t.Run("does not set origin when it doesn't match allowed list", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/festivals", nil)
		req.Header.Set("Origin", "https://evil.com")
		w := httptest.NewRecorder()

		enableCORS(w, req, "http://localhost:5173")

		origin := w.Header().Get("Access-Control-Allow-Origin")
		if origin != "" {
			t.Errorf("Expected no CORS origin, got %s", origin)
		}
	})
}
