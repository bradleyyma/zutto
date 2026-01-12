package mcp

// RankingInput defines the input parameters for the get_anime_ranking tool
type RankingInput struct {
	RankingType string `json:"ranking_type" jsonschema:"Type of ranking (all, tv, movie, ova, ona, special, bypopularity, favorite)"`
	Limit       int    `json:"limit" jsonschema:"Maximum number of results to return (1-100)"`
	Offset      int    `json:"offset" jsonschema:"Offset for pagination"`
}

type DetailsInput struct {
	ID int `json:"id" jsonschema:"ID of the anime to get details for"`
}

type SearchInput struct {
	Query string `json:"query" jsonschema:"Search query for the anime"`
	Limit int    `json:"limit" jsonschema:"Maximum number of results to return (1-50)"`
}
