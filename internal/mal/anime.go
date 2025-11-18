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

// Picture contains image URLs
type Picture struct {
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

// AlternativeTitles contains alternative titles
type AlternativeTitles struct {
	Synonyms []string `json:"synonyms"`
	En       string   `json:"en"`
	Ja       string   `json:"ja"`
}

// Paging contains pagination information
type Paging struct {
	Next string `json:"next,omitempty"`
}

func (a *AnimeService) Search(query string) (*AnimeSearchResponse, error) {
	// URL encode the query parameter
	params := url.Values{}
	params.Add("q", query)
	params.Add("limit", "10") // Optional: limit results

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
