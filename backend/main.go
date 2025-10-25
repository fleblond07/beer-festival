package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config := getConfig()

	festivals, err := loadFestivals(config.FestivalsFile)
	if err != nil {
		log.Fatalf("Failed to load festivals: %v", err)
	}

	log.Printf("Loaded %d festivals from %s", len(festivals), config.FestivalsFile)

	mux := http.NewServeMux()
	mux.HandleFunc(HealthPath, healthCheckHandler)
	mux.HandleFunc(FestivalsPath, makeFestivalsHandler(festivals, config.AllowedOrigins))

	handler := chainMiddleware(mux, requestIDMiddleware, metricsMiddleware, gzipMiddleware)

	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on %s...", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
