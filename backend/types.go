package main

import "time"

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Festival struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
	City         string    `json:"city"`
	Region       string    `json:"region"`
	Location     Location  `json:"location"`
	Image        string    `json:"image"`
	Website      string    `json:"website"`
	BreweryCount int       `json:"breweryCount"`
}

type FestivalDB struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	City        string  `json:"city"`
	Region      string  `json:"region"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Image       string  `json:"image"`
	Website     string  `json:"website"`
}

type FestivalBrewery struct {
	FestivalID int64 `json:"festival_id"`
	BreweryID  int64 `json:"brewery_id"`
}

type Config struct {
	Port           string
	AllowedOrigins string
	SupabaseURL    string
	SupabaseKey    string
}
