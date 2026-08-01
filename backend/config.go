package main

import "os"

func getConfig() Config {
	port := os.Getenv("PORT")
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	databaseURL := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")

	return Config{
		Port:           port,
		AllowedOrigins: allowedOrigins,
		DatabaseURL:    databaseURL,
		JWTSecret:      jwtSecret,
	}
}
