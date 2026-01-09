package main

import (
	"context"
	"log"

	"github.com/davidl71/devwisdom-go/internal/mcp"
)

func main() {
	// Create MCP server using SDK adapter
	server := mcp.NewWisdomServerSDK()

	// Run server with stdio transport (handled by SDK)
	if err := server.Run(context.Background()); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
