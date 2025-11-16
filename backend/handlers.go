package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func enableCORS(w http.ResponseWriter, r *http.Request, allowedOrigins string) {
	origin := r.Header.Get(HeaderOrigin)

	if allowedOrigins == DefaultAllowedOrigins {
		w.Header().Set(HeaderCORSOrigin, DefaultAllowedOrigins)
	} else if origin != "" {
		origins := strings.Split(allowedOrigins, ",")
		for _, allowedOrigin := range origins {
			if strings.TrimSpace(allowedOrigin) == origin {
				w.Header().Set(HeaderCORSOrigin, origin)
				break
			}
		}
	}

	w.Header().Set(HeaderCORSMethods, CORSMethods)
	w.Header().Set(HeaderCORSHeaders, CORSHeaders)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"version": AppVersion,
	})
}

func makeFestivalsHandler(db DatabaseInterface, allowedOrigins string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r, allowedOrigins)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		festivals, err := db.GetFestivals()
		if err != nil {
			log.Printf("Error fetching festivals from database: %v", err)
			http.Error(w, DefaultErrorMessage, http.StatusInternalServerError)
			return
		}

		w.Header().Set(HeaderContentType, ContentTypeJSON)
		if err := json.NewEncoder(w).Encode(festivals); err != nil {
			log.Printf("Error encoding festivals: %v", err)
			http.Error(w, DefaultErrorMessage, http.StatusInternalServerError)
			return
		}
	}
}

func makeLoginHandler(db DatabaseInterface, allowedOrigins string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r, allowedOrigins)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var loginReq LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if loginReq.Email == "" || loginReq.Password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		loginResp, err := db.Login(loginReq.Email, loginReq.Password)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		w.Header().Set(HeaderContentType, ContentTypeJSON)
		if err := json.NewEncoder(w).Encode(loginResp); err != nil {
			log.Printf("Error encoding login response: %v", err)
			http.Error(w, DefaultErrorMessage, http.StatusInternalServerError)
			return
		}
	}
}

func makeVerifyHandler(db DatabaseInterface, allowedOrigins string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r, allowedOrigins)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set(HeaderContentType, ContentTypeJSON)
			json.NewEncoder(w).Encode(VerifyResponse{
				Valid: false,
				Error: "No authorization header",
			})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			w.Header().Set(HeaderContentType, ContentTypeJSON)
			json.NewEncoder(w).Encode(VerifyResponse{
				Valid: false,
				Error: "Invalid authorization format",
			})
			return
		}

		user, err := db.VerifyToken(token)
		if err != nil {
			w.Header().Set(HeaderContentType, ContentTypeJSON)
			json.NewEncoder(w).Encode(VerifyResponse{
				Valid: false,
				Error: "Invalid token",
			})
			return
		}

		w.Header().Set(HeaderContentType, ContentTypeJSON)
		json.NewEncoder(w).Encode(VerifyResponse{
			Valid: true,
			User:  user,
		})
	}
}

func makeBreweriesHandler(db DatabaseInterface, allowedOrigins string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r, allowedOrigins)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		festivalID := strings.TrimPrefix(r.URL.Path, FestivalsBreweriesPath)
		festivalID = strings.TrimSuffix(festivalID, "/breweries")

		if festivalID == "" {
			http.Error(w, "Festival ID is required", http.StatusBadRequest)
			return
		}

		breweries, err := db.GetBreweriesByFestival(festivalID)
		if err != nil {
			log.Printf("Error fetching breweries for festival %s: %v", festivalID, err)
			http.Error(w, DefaultErrorMessage, http.StatusInternalServerError)
			return
		}

		w.Header().Set(HeaderContentType, ContentTypeJSON)
		if err := json.NewEncoder(w).Encode(breweries); err != nil {
			log.Printf("Error encoding breweries: %v", err)
			http.Error(w, DefaultErrorMessage, http.StatusInternalServerError)
			return
		}
	}
}
