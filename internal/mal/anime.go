package mal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type AnimeService struct {
	client *Client
}

// AnimeSearchResponse represents the MAL API response structure
type AnimeSearchResponse struct {
	Data   []AnimeData `json:"data"`
	Paging Paging      `json:"paging"`
}

// AnimeData represents individual anime items in the response
type AnimeData struct {
	Node AnimeNode `json:"node"`
}

// AnimeNode contains the actual anime information
type AnimeNode struct {
	ID                int               `json:"id"`
	Title             string            `json:"title"`
	MainPicture       Picture           `json:"main_picture,omitempty"`
	AlternativeTitles AlternativeTitles `json:"alternative_titles,omitempty"`
}

type AnimeDetails struct {
	Title       string  `json:"title"`
	StartDate   string  `json:"start_date,omitempty"`
	EndDate     string  `json:"end_date,omitempty"`
	Mean        float64 `json:"mean,omitempty"`
	Rank        int     `json:"rank,omitempty"`
	Popularity  int     `json:"popularity,omitempty"`
	Status      string  `json:"status,omitempty"`
	NumEpisodes int     `json:"num_episodes,omitempty"`
	Synopsis    string  `json:"synopsis,omitempty"`
}

// Picture contains image URLs
type Picture struct {
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

// AlternativeTitles contains alternative titles
type AlternativeTitles struct {
	Synonyms *[]string `json:"synonyms,omitempty"`
	En       string    `json:"en,omitempty"`
	Ja       string    `json:"ja,omitempty"`
}

// Paging contains pagination information
type Paging struct {
	Next string `json:"next,omitempty"`
}

type Ranking struct {
	Rank int `json:"rank"`
}

type AnimeRankingData struct {
	Node    AnimeNode `json:"node"`
	Ranking Ranking   `json:"ranking"`
}

type AnimeRankingResponse struct {
	Data   []AnimeRankingData `json:"data"`
	Paging Paging             `json:"paging"`
}

func (a *AnimeService) Search(query string, limit int) (*AnimeSearchResponse, error) {
	// URL encode the query parameter
	params := url.Values{}
	params.Add("q", query)
	params.Add("limit", fmt.Sprintf("%d", limit)) // Optional: limit results

	reqURL := a.client.baseURL.String() + "anime?" + params.Encode()

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check for successful status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read and parse the JSON response
	var searchResponse AnimeSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &searchResponse, nil
}

func (a *AnimeService) Details(animeID int) (*AnimeDetails, error) {
	field := "id,title,synopsis,num_episodes,status,start_date,end_date,mean,rank,popularity"
	reqURL := a.client.baseURL.String() + "anime/" + fmt.Sprintf("%d", animeID) + "?fields=" + field

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var details AnimeDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &details, nil
}

func (a *AnimeService) Rankings(rankingType string, limit, offset int) (*AnimeRankingResponse, error) {
	reqURL := a.client.baseURL.String() + "anime/ranking?ranking_type=" + rankingType + fmt.Sprintf("&limit=%d&offset=%d", limit, offset)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var rankings AnimeRankingResponse
	if err := json.NewDecoder(resp.Body).Decode(&rankings); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &rankings, nil
}

func ValidateAnimeRankingType(rankingType string) error {
	validTypes := map[string]bool{
		"all":          true,
		"tv":           true,
		"movie":        true,
		"ova":          true,
		"ona":          true,
		"special":      true,
		"bypopularity": true,
		"favorite":     true,
	}
	if !validTypes[rankingType] {
		return fmt.Errorf("invalid ranking type: %s", rankingType)
	}
	return nil
}
