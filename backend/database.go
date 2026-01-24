package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/supabase-community/supabase-go"
)

type Database struct {
	client *supabase.Client
	url    string
	key    string
}

func NewDatabase(supabaseURL, supabaseKey string) (*Database, error) {
	if supabaseURL == "" || supabaseKey == "" {
		return nil, fmt.Errorf("SUPABASE_URL and SUPABASE_KEY environment variables are required")
	}

	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Supabase client: %w", err)
	}

	return &Database{
		client: client,
		url:    supabaseURL,
		key:    supabaseKey,
	}, nil
}

type BreweryCount struct {
	FestivalID int64 `json:"festival_id"`
	Count      int64 `json:"count"`
}

type FestivalCount struct {
	BreweryID int64 `json:"brewery_id"`
	Count int64 `json:"count"`
}

func (db *Database) GetFestivals() ([]Festival, error) {
	var festivalsDB []FestivalDB
	_, err := db.client.From("festivals").Select("*", "", false).ExecuteTo(&festivalsDB)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch festivals: %w", err)
	}

	breweryCounts, err := db.getBreweryCounts()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch brewery counts: %w", err)
	}

	festivals := make([]Festival, len(festivalsDB))
	for i, fdb := range festivalsDB {
		startDate, _ := ConvertTime(fdb.StartDate)
		endDate, _ := ConvertTime(fdb.EndDate)

		festivals[i] = Festival{
			ID:          fdb.ID,
			Name:        fdb.Name,
			Description: fdb.Description,
			StartDate:   startDate,
			EndDate:     endDate,
			City:        fdb.City,
			Region:      fdb.Region,
			Location: Location{
				Latitude:  fdb.Latitude,
				Longitude: fdb.Longitude,
			},
			Image:        fdb.Image,
			Website:      fdb.Website,
			BreweryCount: breweryCounts[fdb.ID],
		}
	}

	return festivals, nil
}

func (db *Database) getBreweryCounts() (map[int64]int, error) {
	var counts []BreweryCount
	rpcResult := db.client.Rpc("get_festival_brewery_counts", "", nil)

	err := json.Unmarshal([]byte(rpcResult), &counts)

	if err == nil && len(counts) > 0 {
		result := make(map[int64]int)
		for _, c := range counts {
			result[c.FestivalID] = int(c.Count)
		}
		return result, nil
	}
	return nil, err
}

func (db *Database) getFestivalCounts() (map[int64]int, error) {
	var counts []FestivalCount
	rpcResult := db.client.Rpc("get_brewery_festival_counts", "", nil)

	err := json.Unmarshal([]byte(rpcResult), &counts)

	if err == nil && len(counts) > 0 {
		result := make(map[int64]int)
		for _, c := range counts {
			result[c.BreweryID] = int(c.Count)
		}
		return result, nil
	}
	return nil, err
}
func (db *Database) Login(email, password string) (*LoginResponse, error) {
	resp, err := db.client.Auth.SignInWithEmailPassword(email, password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return &LoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		User: User{
			ID:    resp.User.ID.String(),
			Email: resp.User.Email,
		},
	}, nil
}

func (db *Database) VerifyToken(token string) (*User, error) {
	req, err := http.NewRequest("GET", db.url+"/auth/v1/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("apikey", db.key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token: status %d", resp.StatusCode)
	}

	var userResp struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	return &User{
		ID:    userResp.ID,
		Email: userResp.Email,
	}, nil
}

func (db *Database) GetBreweriesByFestival(festivalID string) ([]Brewery, error) {
	type FestivalBreweryWithBrewery struct {
		BreweryID int64     `json:"brewery_id"`
		Breweries BreweryDB `json:"breweries"`
	}

	var festivalBreweries []FestivalBreweryWithBrewery

	_, err := db.client.From("festivals_breweries").
		Select("brewery_id, breweries(*)", "", false).
		Eq("festival_id", festivalID).
		ExecuteTo(&festivalBreweries)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch breweries: %w", err)
	}

	if len(festivalBreweries) == 0 {
		return []Brewery{}, nil
	}

	breweries := make([]Brewery, len(festivalBreweries))
	for index, brewery := range festivalBreweries {
		breweries[index] = Brewery{
			ID:          brewery.Breweries.ID,
			Name:        brewery.Breweries.Name,
			Description: brewery.Breweries.Description,
			City:        brewery.Breweries.City,
			Website:     brewery.Breweries.Website,
			Logo:        brewery.Breweries.Logo,
		}
	}

	return breweries, nil
}

func (db *Database) GetBreweries() ([]Brewery, error) {
	var breweriesDb []BreweryDB
	_, err := db.client.From("breweries").Select("*", "", false).ExecuteTo(&breweriesDb)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch breweries: %w", err)
	}

	festivalCounts, err := db.getFestivalCounts()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch brewery counts: %w", err)
	}

	breweries := make([]Brewery, len(breweriesDb))
	for i, brewery := range breweriesDb {
		breweries[i] = Brewery{
			ID:          brewery.ID,
			Name:        brewery.Name,
			Description: brewery.Description,
			City:        brewery.City,
			Logo:        brewery.Logo,
			Website:      brewery.Website,
			FestivalCount: festivalCounts[brewery.ID],
		}
	}

	return breweries, nil
}

func (db *Database) CreateFestival(festival *FestivalDB) (*FestivalDB, error) {
	var result []FestivalDB
	_, err := db.client.From("festivals").
		Insert(festival, false, "", "", "").
		ExecuteTo(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to create festival: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no festival returned after creation")
	}

	return &result[0], nil
}
