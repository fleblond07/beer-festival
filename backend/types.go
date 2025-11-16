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

type Brewery struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	City        string `json:"city"`
	Website     string `json:"website"`
	Logo        string `json:"logo"`
}

type BreweryDB struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	City        string `json:"city"`
	Website     string `json:"website"`
	Logo        string `json:"logo"`
}

type Config struct {
	Port           string
	AllowedOrigins string
	SupabaseURL    string
	SupabaseKey    string
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	User         User   `json:"user"`
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type VerifyResponse struct {
	Valid bool   `json:"valid"`
	User  *User  `json:"user,omitempty"`
	Error string `json:"error,omitempty"`
}

type DatabaseInterface interface {
	Login(email, password string) (*LoginResponse, error)
	VerifyToken(token string) (*User, error)
	GetFestivals() ([]Festival, error)
	GetBreweriesByFestival(festivalID string) ([]Brewery, error)
	CreateFestival(festival *FestivalDB) (*FestivalDB, error)
}
