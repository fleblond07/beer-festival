package main

import (
	"encoding/json"
	"fmt"

	"github.com/supabase-community/supabase-go"
)

type Database struct {
	client *supabase.Client
}

func NewDatabase(supabaseURL, supabaseKey string) (*Database, error) {
	if supabaseURL == "" || supabaseKey == "" {
		return nil, fmt.Errorf("SUPABASE_URL and SUPABASE_KEY environment variables are required")
	}

	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Supabase client: %w", err)
	}

	return &Database{client: client}, nil
}

type BreweryCount struct {
	FestivalID int64 `json:"festival_id"`
	Count      int64 `json:"count"`
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
		start_date, _ := ConvertTime(fdb.StartDate)
		end_date, _ := ConvertTime(fdb.EndDate)

		festivals[i] = Festival{
			ID:          fdb.ID,
			Name:        fdb.Name,
			Description: fdb.Description,
			StartDate:   start_date,
			EndDate:     end_date,
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
