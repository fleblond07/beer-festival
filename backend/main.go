package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Festival struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	StartDate    string   `json:"startDate"`
	EndDate      string   `json:"endDate"`
	City         string   `json:"city"`
	Region       string   `json:"region"`
	Location     Location `json:"location"`
	Image        string   `json:"image"`
	Website      string   `json:"website"`
	BreweryCount int      `json:"breweryCount"`
}

type Config struct {
	Port           string
	FestivalsFile  string
	AllowedOrigins string
}

func getConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	festivalsFile := os.Getenv("FESTIVALS_FILE")
	if festivalsFile == "" {
		festivalsFile = "./festivals.json"
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	return Config{
		Port:           port,
		FestivalsFile:  festivalsFile,
		AllowedOrigins: allowedOrigins,
	}
}

func loadFestivals(filePath string) ([]Festival, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read festivals file: %w", err)
	}

	var festivals []Festival
	if err := json.Unmarshal(file, &festivals); err != nil {
		return nil, fmt.Errorf("failed to parse festivals JSON: %w", err)
	}

	return festivals, nil
}

func enableCORS(w http.ResponseWriter, allowedOrigins string) {
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func makeFestivalsHandler(festivals []Festival, allowedOrigins string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, allowedOrigins)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(festivals); err != nil {
			log.Printf("Error encoding festivals: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	config := getConfig()

	festivals, err := loadFestivals(config.FestivalsFile)
	if err != nil {
		log.Fatalf("Failed to load festivals: %v", err)
	}

	log.Printf("Loaded %d festivals from %s", len(festivals), config.FestivalsFile)

	http.HandleFunc("/api/festivals", makeFestivalsHandler(festivals, config.AllowedOrigins))

	addr := ":" + config.Port
	log.Printf("Server starting on %s...", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
