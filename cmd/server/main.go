package main

import (
	"context"
	"log"
	"os"

	"github.com/davidl71/devwisdom-go/internal/mcp"
)

func main() {
	// Create MCP server
	server := mcp.NewWisdomServer()

	// Run server on stdio (standard MCP transport)
	if err := server.Run(context.Background(), os.Stdin, os.Stdout); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
