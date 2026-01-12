/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// mangaCmd represents the manga command
var mangaCmd = &cobra.Command{
	Use:   "manga",
	Short: "Manage and search for manga",
	Long: `Manage and search for manga on MyAnimeList.

Available subcommands:
  search  - Search for manga by query
  ranking - Get manga rankings
  detail  - Get detailed information about a manga

Examples:
  zutto manga search "attack on titan"
  zutto manga ranking --type all
  zutto manga detail --id 12345`,
}

// mangaSearchCmd represents the manga search command
var mangaSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for manga",
	Long: `Search for manga on MyAnimeList.
	
Examples:
  zutto manga search "attack on titan"
  zutto manga search bleach --limit 20`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := strings.Join(args, " ")
		fmt.Printf("Searching for manga: %s\n", query)
		// TODO: Implement manga search functionality
	},
}

func init() {
	rootCmd.AddCommand(mangaCmd)

	// Add subcommands
	mangaCmd.AddCommand(mangaSearchCmd)

	// Search flags
	mangaSearchCmd.Flags().IntP("limit", "l", 10, "Maximum number of results to return (1-50)")
}
