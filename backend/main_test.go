package main

import (
	"bytes"
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
	t.Run("sets CORS headers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/festivals", nil)
		w := httptest.NewRecorder()

		mockDB := &MockDatabase{}

		handler := makeFestivalsHandler(mockDB, "*")
		handler(w, req)

		origin := w.Header().Get("Access-Control-Allow-Origin")
		if origin != "*" {
			t.Errorf("Expected CORS origin *, got %s", origin)
		}

		methods := w.Header().Get("Access-Control-Allow-Methods")
		if methods != "GET, POST, OPTIONS" {
			t.Errorf("Expected CORS methods 'GET, POST, OPTIONS', got %s", methods)
		}
	})

	t.Run("handles OPTIONS requests", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/api/festivals", nil)
		w := httptest.NewRecorder()

		mockDB := &MockDatabase{}

		handler := makeFestivalsHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for OPTIONS, got %d", w.Code)
		}
	})

	t.Run("uses custom allowed origins", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/festivals", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		w := httptest.NewRecorder()

		mockDB := &MockDatabase{}

		handler := makeFestivalsHandler(mockDB, "http://localhost:5173")
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
		mockDB := &MockDatabase{}

		handler := makeFestivalsHandler(mockDB, "http://localhost:5173")
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

		mockDB := &MockDatabase{}

		handler := makeFestivalsHandler(mockDB, "http://localhost:5173,http://localhost:3000")

		handler(w, req)

		origin := w.Header().Get("Access-Control-Allow-Origin")
		if origin != "http://localhost:3000" {
			t.Errorf("Expected CORS origin http://localhost:3000, got %s", origin)
		}
	})

	t.Run("handles GET request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/festivals", nil)
		w := httptest.NewRecorder()
		mockDB := &MockDatabase{}

		handler := makeFestivalsHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
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

type MockDatabase struct {
	loginFunc                  func(email, password string) (*LoginResponse, error)
	verifyTokenFunc            func(token string) (*User, error)
	getFestivalsFunc           func() ([]Festival, error)
	getBreweriesByFestivalFunc func(festivalID string) ([]Brewery, error)
}

func (m *MockDatabase) Login(email, password string) (*LoginResponse, error) {
	if m.loginFunc != nil {
		return m.loginFunc(email, password)
	}
	return nil, nil
}

func (m *MockDatabase) VerifyToken(token string) (*User, error) {
	if m.verifyTokenFunc != nil {
		return m.verifyTokenFunc(token)
	}
	return nil, nil
}

func (m *MockDatabase) GetFestivals() ([]Festival, error) {
	if m.getFestivalsFunc != nil {
		return m.getFestivalsFunc()
	}
	return nil, nil
}

func (m *MockDatabase) GetBreweriesByFestival(festivalID string) ([]Brewery, error) {
	if m.getBreweriesByFestivalFunc != nil {
		return m.getBreweriesByFestivalFunc(festivalID)
	}
	return nil, nil
}

func TestLoginHandler(t *testing.T) {
	t.Run("returns token on successful login", func(t *testing.T) {
		mockDB := &MockDatabase{
			loginFunc: func(email, password string) (*LoginResponse, error) {
				if email == "test@example.com" && password == "password123" {
					return &LoginResponse{
						AccessToken:  "mock-access-token",
						RefreshToken: "mock-refresh-token",
						User: User{
							ID:    "user-123",
							Email: "test@example.com",
						},
					}, nil
				}
				return nil, nil
			},
		}

		body := []byte(`{"email":"test@example.com","password":"password123"}`)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler := makeLoginHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response LoginResponse
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.AccessToken != "mock-access-token" {
			t.Errorf("Expected access token 'mock-access-token', got %s", response.AccessToken)
		}

		if response.User.Email != "test@example.com" {
			t.Errorf("Expected user email 'test@example.com', got %s", response.User.Email)
		}
	})

	t.Run("returns 401 on invalid credentials", func(t *testing.T) {
		mockDB := &MockDatabase{
			loginFunc: func(email, password string) (*LoginResponse, error) {
				return nil, &DatabaseError{Message: "invalid credentials"}
			},
		}

		body := []byte(`{"email":"test@example.com","password":"wrongpassword"}`)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler := makeLoginHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("returns 400 on missing email", func(t *testing.T) {
		mockDB := &MockDatabase{}

		body := []byte(`{"password":"password123"}`)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler := makeLoginHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("returns 400 on missing password", func(t *testing.T) {
		mockDB := &MockDatabase{}

		body := []byte(`{"email":"test@example.com"}`)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler := makeLoginHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("returns 405 on non-POST request", func(t *testing.T) {
		mockDB := &MockDatabase{}

		req := httptest.NewRequest("GET", "/api/auth/login", nil)
		w := httptest.NewRecorder()

		handler := makeLoginHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", w.Code)
		}
	})

	t.Run("handles OPTIONS request", func(t *testing.T) {
		mockDB := &MockDatabase{}

		req := httptest.NewRequest("OPTIONS", "/api/auth/login", nil)
		w := httptest.NewRecorder()

		handler := makeLoginHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for OPTIONS, got %d", w.Code)
		}
	})
}

func TestVerifyHandler(t *testing.T) {
	t.Run("returns valid response with valid token", func(t *testing.T) {
		mockDB := &MockDatabase{
			verifyTokenFunc: func(token string) (*User, error) {
				if token == "valid-token" {
					return &User{
						ID:    "user-123",
						Email: "test@example.com",
					}, nil
				}
				return nil, &DatabaseError{Message: "invalid token"}
			},
		}

		req := httptest.NewRequest("GET", "/api/auth/verify", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()

		handler := makeVerifyHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response VerifyResponse
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if !response.Valid {
			t.Error("Expected valid to be true")
		}

		if response.User == nil {
			t.Fatal("Expected user to be present")
		}

		if response.User.Email != "test@example.com" {
			t.Errorf("Expected user email 'test@example.com', got %s", response.User.Email)
		}
	})

	t.Run("returns invalid response with invalid token", func(t *testing.T) {
		mockDB := &MockDatabase{
			verifyTokenFunc: func(token string) (*User, error) {
				return nil, &DatabaseError{Message: "invalid token"}
			},
		}

		req := httptest.NewRequest("GET", "/api/auth/verify", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()

		handler := makeVerifyHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response VerifyResponse
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Valid {
			t.Error("Expected valid to be false")
		}

		if response.Error != "Invalid token" {
			t.Errorf("Expected error 'Invalid token', got %s", response.Error)
		}
	})

	t.Run("returns invalid response without authorization header", func(t *testing.T) {
		mockDB := &MockDatabase{}

		req := httptest.NewRequest("GET", "/api/auth/verify", nil)
		w := httptest.NewRecorder()

		handler := makeVerifyHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response VerifyResponse
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Valid {
			t.Error("Expected valid to be false")
		}

		if response.Error != "No authorization header" {
			t.Errorf("Expected error 'No authorization header', got %s", response.Error)
		}
	})

	t.Run("returns invalid response without Bearer prefix", func(t *testing.T) {
		mockDB := &MockDatabase{}

		req := httptest.NewRequest("GET", "/api/auth/verify", nil)
		req.Header.Set("Authorization", "invalid-token")
		w := httptest.NewRecorder()

		handler := makeVerifyHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response VerifyResponse
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Valid {
			t.Error("Expected valid to be false")
		}

		if response.Error != "Invalid authorization format" {
			t.Errorf("Expected error 'Invalid authorization format', got %s", response.Error)
		}
	})

	t.Run("returns 405 on non-GET request", func(t *testing.T) {
		mockDB := &MockDatabase{}

		req := httptest.NewRequest("POST", "/api/auth/verify", nil)
		w := httptest.NewRecorder()

		handler := makeVerifyHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", w.Code)
		}
	})

	t.Run("handles OPTIONS request", func(t *testing.T) {
		mockDB := &MockDatabase{}

		req := httptest.NewRequest("OPTIONS", "/api/auth/verify", nil)
		w := httptest.NewRecorder()

		handler := makeVerifyHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for OPTIONS, got %d", w.Code)
		}
	})
}

type DatabaseError struct {
	Message string
}

func (e *DatabaseError) Error() string {
	return e.Message
}

func TestBreweriesHandler(t *testing.T) {
	t.Run("returns breweries for a festival", func(t *testing.T) {
		mockDB := &MockDatabase{
			getBreweriesByFestivalFunc: func(festivalID string) ([]Brewery, error) {
				if festivalID == "1" {
					return []Brewery{
						{
							ID:          1,
							Name:        "Test Brewery",
							Description: "A test brewery",
							City:        "Paris",
							Website:     "https://example.com",
							Logo:        "https://example.com/logo.png",
						},
					}, nil
				}
				return nil, &DatabaseError{Message: "festival not found"}
			},
		}

		req := httptest.NewRequest("GET", "/api/festivals/1/breweries", nil)
		w := httptest.NewRecorder()

		handler := makeBreweriesHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var breweries []Brewery
		if err := json.NewDecoder(w.Body).Decode(&breweries); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(breweries) != 1 {
			t.Errorf("Expected 1 brewery, got %d", len(breweries))
		}

		if breweries[0].Name != "Test Brewery" {
			t.Errorf("Expected brewery name 'Test Brewery', got %s", breweries[0].Name)
		}
	})

	t.Run("returns error for invalid festival ID", func(t *testing.T) {
		mockDB := &MockDatabase{
			getBreweriesByFestivalFunc: func(festivalID string) ([]Brewery, error) {
				return nil, &DatabaseError{Message: "festival not found"}
			},
		}

		req := httptest.NewRequest("GET", "/api/festivals/999/breweries", nil)
		w := httptest.NewRecorder()

		handler := makeBreweriesHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
	})

	t.Run("returns 400 for missing festival ID", func(t *testing.T) {
		mockDB := &MockDatabase{}

		req := httptest.NewRequest("GET", "/api/festivals//breweries", nil)
		w := httptest.NewRecorder()

		handler := makeBreweriesHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("returns 405 on non-GET request", func(t *testing.T) {
		mockDB := &MockDatabase{}

		req := httptest.NewRequest("POST", "/api/festivals/1/breweries", nil)
		w := httptest.NewRecorder()

		handler := makeBreweriesHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", w.Code)
		}
	})

	t.Run("handles OPTIONS request", func(t *testing.T) {
		mockDB := &MockDatabase{}

		req := httptest.NewRequest("OPTIONS", "/api/festivals/1/breweries", nil)
		w := httptest.NewRecorder()

		handler := makeBreweriesHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for OPTIONS, got %d", w.Code)
		}
	})
}

func TestFestivalsHandlerWithDB(t *testing.T) {
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

	t.Run("returns festivals from database", func(t *testing.T) {
		mockDB := &MockDatabase{
			getFestivalsFunc: func() ([]Festival, error) {
				return mockFestivals, nil
			},
		}

		req := httptest.NewRequest("GET", "/api/festivals", nil)
		w := httptest.NewRecorder()

		handler := makeFestivalsHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var festivals []Festival
		if err := json.NewDecoder(w.Body).Decode(&festivals); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(festivals) != 1 {
			t.Errorf("Expected 1 festival, got %d", len(festivals))
		}

		if festivals[0].Name != "Test Festival" {
			t.Errorf("Expected festival name 'Test Festival', got %s", festivals[0].Name)
		}
	})

	t.Run("returns error when database fails", func(t *testing.T) {
		mockDB := &MockDatabase{
			getFestivalsFunc: func() ([]Festival, error) {
				return nil, &DatabaseError{Message: "database error"}
			},
		}

		req := httptest.NewRequest("GET", "/api/festivals", nil)
		w := httptest.NewRecorder()

		handler := makeFestivalsHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
	})

	t.Run("handles OPTIONS request", func(t *testing.T) {
		mockDB := &MockDatabase{}

		req := httptest.NewRequest("OPTIONS", "/api/festivals", nil)
		w := httptest.NewRecorder()

		handler := makeFestivalsHandler(mockDB, "*")
		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for OPTIONS, got %d", w.Code)
		}
	})
}
