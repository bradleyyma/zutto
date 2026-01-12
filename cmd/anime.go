/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/bradleyyma/zutto/internal/mal"
	"github.com/spf13/cobra"
)

// animeCmd represents the anime command
var animeCmd = &cobra.Command{
	Use:   "anime",
	Short: "Manage and search for anime",
	Long: `Manage and search for anime on MyAnimeList.

Available subcommands:
  search  - Search for anime by query
  ranking - Get anime rankings
  detail  - Get detailed information about an anime

Examples:
  zutto anime search "one piece"
  zutto anime ranking --type tv
  zutto anime detail --id 5114`,
}

// animeSearchCmd represents the anime search command
var animeSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for anime",
	Long: `Search for anime on MyAnimeList.
	
Examples:
  zutto anime search "one piece"
  zutto anime search naruto --limit 20`,
	Args: cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		if limit <= 0 || limit > 50 {
			return fmt.Errorf("limit must be greater than 0 and less than or equal to 50, got %d", limit)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		query := strings.Join(args, " ")
		limit, _ := cmd.Flags().GetInt("limit")

		client := mal.NewClient(nil, "")
		results, err := client.Anime.Search(query, limit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error searching for anime: %v\n", err)
			os.Exit(1)
		}

		if len(results.Data) == 0 {
			fmt.Println("No anime found")
			return
		}

		fmt.Printf("Found %d anime:\n\n", len(results.Data))
		for i, anime := range results.Data {
			fmt.Printf("%d. %s (ID: %d)\n", i+1, anime.Node.Title, anime.Node.ID)
			if anime.Node.AlternativeTitles.En != "" {
				fmt.Printf("   English: %s\n", anime.Node.AlternativeTitles.En)
			}
		}
	},
}

// animeRankingCmd represents the anime ranking command
var animeRankingCmd = &cobra.Command{
	Use:   "ranking",
	Short: "Get anime rankings",
	Long: `Get anime rankings from MyAnimeList.

Examples:
  zutto anime ranking
  zutto anime ranking --type tv
  zutto anime ranking --type movie --limit 20`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		rankingType, _ := cmd.Flags().GetString("type")
		if err := mal.ValidateAnimeRankingType(rankingType); err != nil {
			return err
		}

		limit, _ := cmd.Flags().GetInt("limit")
		if limit <= 0 || limit > 100 {
			return fmt.Errorf("limit must be between 1 and 100, got %d", limit)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		rankingType, _ := cmd.Flags().GetString("type")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		client := mal.NewClient(nil, "")
		rankings, err := client.Anime.Rankings(rankingType, limit, offset)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error retrieving anime rankings: %v\n", err)
			os.Exit(1)
		}

		if len(rankings.Data) == 0 {
			fmt.Println("No anime rankings found")
			return
		}

		fmt.Printf("Top %d Anime Rankings (%s):\n\n", len(rankings.Data), rankingType)
		for _, entry := range rankings.Data {
			fmt.Printf("%d. %s (ID: %d)\n", entry.Ranking.Rank, entry.Node.Title, entry.Node.ID)
			if entry.Node.AlternativeTitles.En != "" {
				fmt.Printf("    English: %s\n", entry.Node.AlternativeTitles.En)
			}
		}
	},
}

// animeDetailCmd represents the anime detail command
var animeDetailCmd = &cobra.Command{
	Use:   "detail",
	Short: "Get detailed information about an anime",
	Long: `Get detailed information about an anime on MyAnimeList by ID or name.

You must provide either --id or --name (but not both).

Examples:
  zutto anime detail --id 5114
  zutto anime detail -i 5114
  zutto anime detail --name "Fullmetal Alchemist: Brotherhood"
  zutto anime detail -n "naruto"`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetInt("id")
		name, _ := cmd.Flags().GetString("name")

		if id == 0 && name == "" {
			return fmt.Errorf("either --id or --name must be provided")
		}
		if id != 0 && name != "" {
			return fmt.Errorf("cannot use both --id and --name flags together")
		}
		if id < 0 {
			return fmt.Errorf("id must be a positive integer, got %d", id)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetInt("id")
		name, _ := cmd.Flags().GetString("name")

		client := mal.NewClient(nil, "")

		// If name is provided, search first to get the ID
		if name != "" {
			fmt.Printf("Searching for anime: %s\n", name)
			searchResults, err := client.Anime.Search(name, 1)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error searching for anime: %v\n", err)
				os.Exit(1)
			}
			if len(searchResults.Data) == 0 {
				fmt.Fprintf(os.Stderr, "No anime found with name: %s\n", name)
				os.Exit(1)
			}
			id = searchResults.Data[0].Node.ID
			fmt.Printf("Found: %s (ID: %d)\n\n", searchResults.Data[0].Node.Title, id)
		}

		// Get anime details
		detail, err := client.Anime.Details(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting anime details: %v\n", err)
			os.Exit(1)
		}

		// Display results
		fmt.Printf("Title: %s\n", detail.Title)
		if detail.Synopsis != "" {
			fmt.Printf("Synopsis: %s\n", detail.Synopsis)
		}
		fmt.Printf("Number of Episodes: %d\n", detail.NumEpisodes)
		fmt.Printf("Status: %s\n", detail.Status)
		if detail.Mean > 0 {
			fmt.Printf("Score: %.2f\n", detail.Mean)
		}
		if detail.StartDate != "" {
			fmt.Printf("Start Date: %s\n", detail.StartDate)
		}
		if detail.EndDate != "" {
			fmt.Printf("End Date: %s\n", detail.EndDate)
		}
	},
}

func init() {
	rootCmd.AddCommand(animeCmd)

	// Add subcommands
	animeCmd.AddCommand(animeSearchCmd)
	animeCmd.AddCommand(animeRankingCmd)
	animeCmd.AddCommand(animeDetailCmd)

	// Search flags
	animeSearchCmd.Flags().IntP("limit", "l", 10, "Maximum number of results to return (1-50)")

	// Ranking flags
	animeRankingCmd.Flags().String("type", "all", "Type of ranking (all, tv, movie, ova, ona, special, bypopularity, favorite)")
	animeRankingCmd.Flags().IntP("limit", "l", 50, "Maximum number of results to return (1-100)")
	animeRankingCmd.Flags().Int("offset", 0, "Offset for pagination")

	// Detail flags
	animeDetailCmd.Flags().IntP("id", "i", 0, "Anime ID")
	animeDetailCmd.Flags().StringP("name", "n", "", "Anime name (will search and use first result)")
}
