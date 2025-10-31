package main

import "os"

func getConfig() Config {
	port := os.Getenv("PORT")
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	return Config{
		Port:           port,
		AllowedOrigins: allowedOrigins,
		SupabaseURL:    supabaseURL,
		SupabaseKey:    supabaseKey,
	}
}
