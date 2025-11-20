/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/bradleyyma/zutto/internal/mal"
	"github.com/spf13/cobra"
)

// detailAnimeCmd represents the detail anime command
var detailAnimeCmd = &cobra.Command{
	Use:   "anime",
	Short: "Get details for anime",
	Long: `Get details for an anime on MyAnimeList by ID or name.

You must provide either --id or --name (but not both).

Examples:
  zutto detail anime --id 5114
  zutto detail anime -i 5114
  zutto detail anime --name "Fullmetal Alchemist: Brotherhood"
  zutto detail anime -n "naruto"`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetInt("id")
		name, _ := cmd.Flags().GetString("name")

		// Check that exactly one flag is provided
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

		// Get MAL client ID from environment variable
		clientID := os.Getenv("MAL_CLIENT_ID")
		client := mal.NewClient(nil, clientID)

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
		fmt.Printf("Synopsis: %s\n", detail.Synopsis)
		fmt.Printf("Number of Episodes: %d\n", detail.NumEpisodes)
		fmt.Printf("Status: %s\n", detail.Status)
		fmt.Printf("Score: %.2f\n", detail.Mean)
		fmt.Printf("Start Date: %s\n", detail.StartDate)
		fmt.Printf("End Date: %s\n", detail.EndDate)
	},
}

func init() {
	detailCmd.AddCommand(detailAnimeCmd)

	// Add mutually exclusive flags
	detailAnimeCmd.Flags().IntP("id", "i", 0, "Anime ID")
	detailAnimeCmd.Flags().StringP("name", "n", "", "Anime name (will search and use first result)")
}
