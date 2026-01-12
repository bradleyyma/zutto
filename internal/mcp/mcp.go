package mcp

import (
	"context"
	"fmt"

	"github.com/bradleyyma/zutto/internal/mal"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server wraps the MCP server with MAL client
type Server struct {
	mcpServer *mcp.Server
	malClient *mal.Client
}

// NewMCPServer creates and configures the MCP server with all tools
func NewMCPServer() (*Server, error) {
	// Initialize MAL client - it will automatically get clientID from MAL_CLIENT_ID env var
	malClient := mal.NewClient(nil, "")

	// Validate that we have a client ID
	if malClient == nil {
		return nil, fmt.Errorf("failed to create MAL client")
	}

	// Create MCP server
	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    "zutto",
			Version: "1.0.0",
		},
		nil,
	)

	server := &Server{
		mcpServer: mcpServer,
		malClient: malClient,
	}

	// Register tools
	if err := server.registerTools(); err != nil {
		return nil, fmt.Errorf("failed to register tools: %w", err)
	}

	return server, nil
}

// registerTools registers all available MCP tools
func (s *Server) registerTools() error {
	// Register the anime ranking tool using the generic AddTool
	mcp.AddTool(
		s.mcpServer,
		&mcp.Tool{
			Name:        "get_anime_ranking",
			Description: "Get anime rankings from MyAnimeList. Returns the top-ranked anime based on the specified ranking type.",
		},
		s.handleGetAnimeRanking,
	)

	mcp.AddTool(
		s.mcpServer,
		&mcp.Tool{
			Name:        "get_anime_details",
			Description: "Get detailed information about an anime by its MyAnimelistID",
		},
		s.handleAnimeDetails,
	)

	mcp.AddTool(
		s.mcpServer,
		&mcp.Tool{
			Name:        "search_anime",
			Description: "Search for anime on MyAnimeList by query string. Use this tool to find anime ids",
		},
		s.handleAnimeSearch,
	)

	return nil
}

// handleGetAnimeRanking handles the get_anime_ranking tool invocation
func (s *Server) handleGetAnimeRanking(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input RankingInput,
) (*mcp.CallToolResult, mal.AnimeRankingResponse, error) {
	// Apply defaults
	if input.RankingType == "" {
		input.RankingType = "all"
	}
	if input.Limit == 0 {
		input.Limit = 10
	}
	if input.Limit < 1 || input.Limit > 100 {
		return nil, mal.AnimeRankingResponse{}, fmt.Errorf("limit must be between 1 and 100")
	}
	if input.Offset < 0 {
		input.Offset = 0
	}

	// Validate ranking type
	if err := mal.ValidateAnimeRankingType(input.RankingType); err != nil {
		return nil, mal.AnimeRankingResponse{}, err
	}

	// Call MAL API
	rankings, err := s.malClient.Anime.Rankings(input.RankingType, input.Limit, input.Offset)
	if err != nil {
		return nil, mal.AnimeRankingResponse{}, fmt.Errorf("failed to fetch rankings: %w", err)
	}

	return nil, *rankings, nil
}

func (s *Server) handleAnimeDetails(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input DetailsInput,
) (*mcp.CallToolResult, mal.AnimeDetails, error) {
	if input.ID <= 0 {
		return nil, mal.AnimeDetails{}, fmt.Errorf("invalid anime ID")
	}

	details, err := s.malClient.Anime.Details(input.ID)
	if err != nil {
		return nil, mal.AnimeDetails{}, fmt.Errorf("failed to fetch anime details: %w", err)
	}

	return nil, *details, nil
}

func (s *Server) handleAnimeSearch(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SearchInput,
) (*mcp.CallToolResult, mal.AnimeSearchResponse, error) {
	// Apply defaults
	if input.Limit == 0 {
		input.Limit = 10
	}
	if input.Limit < 1 || input.Limit > 50 {
		return nil, mal.AnimeSearchResponse{}, fmt.Errorf("limit must be between 1 and 50")
	}
	if input.Query == "" {
		return nil, mal.AnimeSearchResponse{}, fmt.Errorf("query cannot be empty")
	}

	results, err := s.malClient.Anime.Search(input.Query, input.Limit)
	if err != nil {
		return nil, mal.AnimeSearchResponse{}, fmt.Errorf("failed to fetch anime search results: %w", err)
	}

	return nil, *results, nil
}

// Run starts the MCP server
func (s *Server) Run(ctx context.Context) error {
	return s.mcpServer.Run(ctx, &mcp.StdioTransport{})
}
