package main

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
