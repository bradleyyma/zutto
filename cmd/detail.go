/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var detailCmd = &cobra.Command{
	Use:   "detail",
	Short: "Get details for anime or manga",
	Long: `Get details for anime or manga on MyAnimeList.

Available subcommands:
  anime - Get details for anime
  manga - Get details for manga
Examples:
  zutto detail anime 13469
  zutto detail manga 12345`,
}

func init() {
	rootCmd.AddCommand(detailCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
