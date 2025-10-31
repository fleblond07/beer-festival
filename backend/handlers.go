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

func makeFestivalsHandler(festivals []Festival, allowedOrigins string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r, allowedOrigins)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.Header().Set(HeaderContentType, ContentTypeJSON)
		if err := json.NewEncoder(w).Encode(festivals); err != nil {
			log.Printf("Error encoding festivals: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func makeFestivalsHandlerWithDB(db *Database, allowedOrigins string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r, allowedOrigins)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		festivals, err := db.GetFestivals()
		if err != nil {
			log.Printf("Error fetching festivals from database: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set(HeaderContentType, ContentTypeJSON)
		if err := json.NewEncoder(w).Encode(festivals); err != nil {
			log.Printf("Error encoding festivals: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
