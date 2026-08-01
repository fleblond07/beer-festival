package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	db        *sql.DB
	jwtSecret []byte
}

func NewDatabase(databaseURL, jwtSecret string) (*Database, error) {
	if databaseURL == "" || jwtSecret == "" {
		return nil, fmt.Errorf("DATABASE_URL and JWT_SECRET environment variables are required")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open Postgres connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	return &Database{db: db, jwtSecret: []byte(jwtSecret)}, nil
}

func (db *Database) Close() error {
	return db.db.Close()
}

type BreweryCount struct {
	FestivalID int64 `json:"festival_id"`
	Count      int64 `json:"count"`
}

type FestivalCount struct {
	BreweryID int64 `json:"brewery_id"`
	Count     int64 `json:"count"`
}

func (db *Database) GetFestivals() ([]Festival, error) {
	rows, err := db.db.Query(`
		SELECT
			f.id,
			COALESCE(f.name, ''),
			COALESCE(f.description, ''),
			f.start_date,
			f.end_date,
			COALESCE(f.city, ''),
			COALESCE(f.region, ''),
			COALESCE(f.latitude, 0),
			COALESCE(f.longitude, 0),
			COALESCE(f.image, ''),
			COALESCE(f.website, ''),
			COUNT(fb.brewery_id)
		FROM festivals f
		LEFT JOIN festivals_breweries fb ON fb.festival_id = f.id
		GROUP BY f.id
		ORDER BY f.start_date ASC, f.id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch festivals: %w", err)
	}
	defer rows.Close()

	festivals := []Festival{}
	for rows.Next() {
		var festival Festival
		var startDate time.Time
		var endDate time.Time

		if err := rows.Scan(
			&festival.ID,
			&festival.Name,
			&festival.Description,
			&startDate,
			&endDate,
			&festival.City,
			&festival.Region,
			&festival.Location.Latitude,
			&festival.Location.Longitude,
			&festival.Image,
			&festival.Website,
			&festival.BreweryCount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan festival: %w", err)
		}

		festival.StartDate = startDate
		festival.EndDate = endDate
		festivals = append(festivals, festival)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate festivals: %w", err)
	}

	return festivals, nil
}

func (db *Database) GetBreweriesByFestival(festivalID string) ([]Brewery, error) {
	rows, err := db.db.Query(`
		SELECT
			b.id,
			COALESCE(b.name, ''),
			COALESCE(b.description, ''),
			COALESCE(b.city, ''),
			COALESCE(b.website, ''),
			COALESCE(b.logo, '')
		FROM breweries b
		INNER JOIN festivals_breweries fb ON fb.brewery_id = b.id
		WHERE fb.festival_id = $1
		ORDER BY b.name ASC, b.id ASC
	`, festivalID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch breweries: %w", err)
	}
	defer rows.Close()

	breweries := []Brewery{}
	for rows.Next() {
		var brewery Brewery
		if err := rows.Scan(
			&brewery.ID,
			&brewery.Name,
			&brewery.Description,
			&brewery.City,
			&brewery.Website,
			&brewery.Logo,
		); err != nil {
			return nil, fmt.Errorf("failed to scan brewery: %w", err)
		}
		breweries = append(breweries, brewery)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate breweries: %w", err)
	}

	return breweries, nil
}

func (db *Database) GetBreweries() ([]Brewery, error) {
	rows, err := db.db.Query(`
		SELECT
			b.id,
			COALESCE(b.name, ''),
			COALESCE(b.description, ''),
			COALESCE(b.city, ''),
			COALESCE(b.website, ''),
			COALESCE(b.logo, ''),
			COUNT(fb.festival_id)
		FROM breweries b
		LEFT JOIN festivals_breweries fb ON fb.brewery_id = b.id
		GROUP BY b.id
		ORDER BY b.name ASC, b.id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch breweries: %w", err)
	}
	defer rows.Close()

	breweries := []Brewery{}
	for rows.Next() {
		var brewery Brewery
		if err := rows.Scan(
			&brewery.ID,
			&brewery.Name,
			&brewery.Description,
			&brewery.City,
			&brewery.Website,
			&brewery.Logo,
			&brewery.FestivalCount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan brewery: %w", err)
		}
		breweries = append(breweries, brewery)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate breweries: %w", err)
	}

	return breweries, nil
}

func (db *Database) CreateFestival(festival *FestivalDB) (*FestivalDB, error) {
	var created FestivalDB
	var startDate time.Time
	var endDate time.Time

	err := db.db.QueryRow(`
		INSERT INTO festivals (
			name,
			description,
			start_date,
			end_date,
			city,
			region,
			latitude,
			longitude,
			image,
			website
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, name, description, start_date, end_date, city, region, latitude, longitude, image, website
	`,
		festival.Name,
		festival.Description,
		festival.StartDate,
		festival.EndDate,
		festival.City,
		festival.Region,
		festival.Latitude,
		festival.Longitude,
		festival.Image,
		festival.Website,
	).Scan(
		&created.ID,
		&created.Name,
		&created.Description,
		&startDate,
		&endDate,
		&created.City,
		&created.Region,
		&created.Latitude,
		&created.Longitude,
		&created.Image,
		&created.Website,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create festival: %w", err)
	}

	created.StartDate = startDate.Format(DefaultTimeFormat)
	created.EndDate = endDate.Format(DefaultTimeFormat)

	return &created, nil
}

func (db *Database) Login(email, password string) (*LoginResponse, error) {
	var user User
	err := db.db.QueryRow(`
		SELECT id::text, email
		FROM users
		WHERE lower(email) = lower($1)
		AND password_hash = crypt($2, password_hash)
	`, email, password).Scan(&user.ID, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("authentication failed")
		}
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	accessToken, err := db.signToken(user, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := db.signToken(user, 30*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (db *Database) VerifyToken(token string) (*User, error) {
	payload, err := db.verifySignedToken(token)
	if err != nil {
		return nil, err
	}

	return &User{ID: payload.Subject, Email: payload.Email}, nil
}

type tokenPayload struct {
	Subject   string `json:"sub"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
}

func (db *Database) signToken(user User, ttl time.Duration) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	payload := tokenPayload{
		Subject:   user.ID,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(ttl).Unix(),
		IssuedAt:  time.Now().Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to encode token header: %w", err)
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to encode token payload: %w", err)
	}

	unsigned := base64.RawURLEncoding.EncodeToString(headerJSON) + "." + base64.RawURLEncoding.EncodeToString(payloadJSON)
	signature := db.sign(unsigned)

	return unsigned + "." + signature, nil
}

func (db *Database) verifySignedToken(token string) (*tokenPayload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	unsigned := parts[0] + "." + parts[1]
	expectedSignature := db.sign(unsigned)
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSignature)) {
		return nil, fmt.Errorf("invalid token signature")
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid token payload: %w", err)
	}

	var payload tokenPayload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return nil, fmt.Errorf("failed to decode token payload: %w", err)
	}

	if payload.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("token expired")
	}

	return &payload, nil
}

func (db *Database) sign(unsigned string) string {
	mac := hmac.New(sha256.New, db.jwtSecret)
	mac.Write([]byte(unsigned))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
