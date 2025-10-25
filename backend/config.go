package main

import "os"

func getConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	festivalsFile := os.Getenv("FESTIVALS_FILE")
	if festivalsFile == "" {
		festivalsFile = DefaultFestivalsFile
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = DefaultAllowedOrigins
	}

	return Config{
		Port:           port,
		FestivalsFile:  festivalsFile,
		AllowedOrigins: allowedOrigins,
	}
}
