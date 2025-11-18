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

// searchAnimeCmd represents the search anime command
var searchAnimeCmd = &cobra.Command{
	Use:   "anime [query]",
	Short: "Search for anime",
	Long: `Search for anime on MyAnimeList.
	
Example:
  zutto search anime "one piece"
  zutto search anime naruto`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := strings.Join(args, " ")

		// Get MAL client ID from environment variable
		clientID := os.Getenv("MAL_CLIENT_ID")
		client := mal.NewClient(nil, clientID)

		// Search for anime
		results, err := client.Anime.Search(query)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error searching for anime: %v\n", err)
			os.Exit(1)
		}

		// Display results
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

func init() {
	searchCmd.AddCommand(searchAnimeCmd)

	// Add flags specific to anime search if needed
	// searchAnimeCmd.Flags().StringP("status", "s", "", "Filter by status (watching, completed, etc.)")
}
