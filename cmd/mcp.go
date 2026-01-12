package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/bradleyyma/zutto/internal/mcp"
	"github.com/spf13/cobra"
)

// mcpCmd represents the mcp command
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start the MCP server",
	Long: `Start a Model Context Protocol (MCP) server for Zutto.
	
This command initializes and runs the MCP server, allowing interaction
with the Zutto application via the MCP protocol.

The server provides the following tools:
  - get_anime_ranking: Get anime rankings from MyAnimeList

Example:
  zutto mcp`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create MCP server
		server, err := mcp.NewMCPServer()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating MCP server: %v\n", err)
			os.Exit(1)
		}

		// Run the server
		if err := server.Run(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "Error running MCP server: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
