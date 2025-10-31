package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestIDMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value("requestID")
		if requestID == nil {
			t.Error("Expected requestID in context, got nil")
		}
		w.WriteHeader(http.StatusOK)
	})

	t.Run("generates request ID when not provided", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		middleware := requestIDMiddleware(handler)
		middleware.ServeHTTP(w, req)

		requestID := w.Header().Get("X-Request-ID")
		if requestID == "" {
			t.Error("Expected X-Request-ID header to be set")
		}

		if len(requestID) != 32 {
			t.Errorf("Expected request ID length 32, got %d", len(requestID))
		}
	})

	t.Run("uses provided request ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Request-ID", "custom-id-123")
		w := httptest.NewRecorder()

		middleware := requestIDMiddleware(handler)
		middleware.ServeHTTP(w, req)

		requestID := w.Header().Get("X-Request-ID")
		if requestID != "custom-id-123" {
			t.Errorf("Expected request ID custom-id-123, got %s", requestID)
		}
	})

	t.Run("adds request ID to context", func(t *testing.T) {
		contextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Context().Value("requestID")
			if requestID == nil {
				t.Error("Expected requestID in context")
				return
			}
			if requestID.(string) == "" {
				t.Error("Expected non-empty requestID in context")
			}
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		middleware := requestIDMiddleware(contextHandler)
		middleware.ServeHTTP(w, req)
	})
}

func TestMetricsMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	t.Run("calls next handler", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req = req.WithContext(context.WithValue(req.Context(), "requestID", "test-id"))
		w := httptest.NewRecorder()

		middleware := metricsMiddleware(handler)
		middleware.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		if body != "OK" {
			t.Errorf("Expected body 'OK', got %s", body)
		}
	})

	t.Run("logs request with different status codes", func(t *testing.T) {
		statusCodes := []int{200, 404, 500}

		for _, statusCode := range statusCodes {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(statusCode)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			req = req.WithContext(context.WithValue(req.Context(), "requestID", "test-id"))
			w := httptest.NewRecorder()

			middleware := metricsMiddleware(handler)
			middleware.ServeHTTP(w, req)

			if w.Code != statusCode {
				t.Errorf("Expected status %d, got %d", statusCode, w.Code)
			}
		}
	})
}

func TestGzipMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("This is a test response that should be compressed"))
	})

	t.Run("compresses response when gzip is accepted", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		w := httptest.NewRecorder()

		middleware := gzipMiddleware(handler)
		middleware.ServeHTTP(w, req)

		contentEncoding := w.Header().Get("Content-Encoding")
		if contentEncoding != "gzip" {
			t.Errorf("Expected Content-Encoding gzip, got %s", contentEncoding)
		}

		reader, err := gzip.NewReader(w.Body)
		if err != nil {
			t.Fatalf("Failed to create gzip reader: %v", err)
		}
		defer reader.Close()

		decompressed := &bytes.Buffer{}
		if _, err := io.Copy(decompressed, reader); err != nil {
			t.Fatalf("Failed to decompress response: %v", err)
		}

		expected := "This is a test response that should be compressed"
		if decompressed.String() != expected {
			t.Errorf("Expected decompressed body '%s', got '%s'", expected, decompressed.String())
		}
	})

	t.Run("does not compress when gzip is not accepted", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		middleware := gzipMiddleware(handler)
		middleware.ServeHTTP(w, req)

		contentEncoding := w.Header().Get("Content-Encoding")
		if contentEncoding == "gzip" {
			t.Error("Should not set Content-Encoding to gzip when not accepted")
		}

		body := w.Body.String()
		expected := "This is a test response that should be compressed"
		if body != expected {
			t.Errorf("Expected body '%s', got '%s'", expected, body)
		}
	})
}

func TestChainMiddleware(t *testing.T) {
	t.Run("chains multiple middleware in correct order", func(t *testing.T) {
		var order []string

		middleware1 := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "middleware1-before")
				next.ServeHTTP(w, r)
				order = append(order, "middleware1-after")
			})
		}

		middleware2 := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "middleware2-before")
				next.ServeHTTP(w, r)
				order = append(order, "middleware2-after")
			})
		}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "handler")
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		chained := chainMiddleware(handler, middleware1, middleware2)
		chained.ServeHTTP(w, req)

		expected := []string{
			"middleware1-before",
			"middleware2-before",
			"handler",
			"middleware2-after",
			"middleware1-after",
		}

		if len(order) != len(expected) {
			t.Fatalf("Expected %d calls, got %d", len(expected), len(order))
		}

		for i, v := range expected {
			if order[i] != v {
				t.Errorf("Expected order[%d] to be '%s', got '%s'", i, v, order[i])
			}
		}
	})
}

func TestResponseWriter(t *testing.T) {
	t.Run("captures status code", func(t *testing.T) {
		w := httptest.NewRecorder()
		rw := &ResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		rw.WriteHeader(http.StatusNotFound)

		if rw.statusCode != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", rw.statusCode)
		}
	})

	t.Run("prevents multiple WriteHeader calls", func(t *testing.T) {
		w := httptest.NewRecorder()
		rw := &ResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		rw.WriteHeader(http.StatusNotFound)
		rw.WriteHeader(http.StatusInternalServerError)

		if rw.statusCode != http.StatusNotFound {
			t.Errorf("Expected status code to remain 404, got %d", rw.statusCode)
		}
	})

	t.Run("Write sets status code to 200 if not set", func(t *testing.T) {
		w := httptest.NewRecorder()
		rw := &ResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		rw.Write([]byte("test"))

		if !rw.written {
			t.Error("Expected written to be true")
		}
	})
}

func TestGenerateRequestID(t *testing.T) {
	t.Run("generates unique IDs", func(t *testing.T) {
		id1 := generateRequestID()
		id2 := generateRequestID()

		if id1 == id2 {
			t.Error("Expected unique request IDs, got duplicates")
		}
	})

	t.Run("generates IDs of correct length", func(t *testing.T) {
		id := generateRequestID()

		if len(id) != 32 {
			t.Errorf("Expected request ID length 32, got %d", len(id))
		}
	})

	t.Run("generates hex-encoded IDs", func(t *testing.T) {
		id := generateRequestID()

		for _, c := range id {
			if !strings.ContainsRune("0123456789abcdef", c) {
				t.Errorf("Expected hex character, got %c", c)
			}
		}
	})
}
