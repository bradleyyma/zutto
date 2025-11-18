/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// searchMangaCmd represents the search manga command
var searchMangaCmd = &cobra.Command{
	Use:   "manga [query]",
	Short: "Search for manga",
	Long: `Search for manga on MyAnimeList.
	
Example:
  zutto search manga "attack on titan"
  zutto search manga bleach`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := strings.Join(args, " ")
		fmt.Printf("Searching for manga: %s\n", query)
		// TODO: Implement actual manga search functionality
	},
}

func init() {
	searchCmd.AddCommand(searchMangaCmd)

	// Add flags specific to manga search if needed
	// searchMangaCmd.Flags().StringP("status", "s", "", "Filter by status (reading, completed, etc.)")
}
