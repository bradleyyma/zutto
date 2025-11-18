package mal

import (
	"net/http"
	"net/url"
	"os"
)

const malURL = "https://api.myanimelist.net/v2/"

type Client struct {
	client   *http.Client
	baseURL  *url.URL
	clientID string

	Anime *AnimeService
}

func NewClient(httpClient *http.Client, clientID string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	// If no clientID provided, try to get from environment variable
	if clientID == "" {
		clientID = os.Getenv("MAL_CLIENT_ID")
	}

	baseURL, _ := url.Parse(malURL)
	c := &Client{
		client:   httpClient,
		baseURL:  baseURL,
		clientID: clientID,
	}
	c.Anime = &AnimeService{client: c}

	return c
}

// Do sends an HTTP request and adds the MAL client ID header
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-MAL-CLIENT-ID", c.clientID)
	return c.client.Do(req)
}
