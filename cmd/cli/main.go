package main

import (
	"context"
	"fmt"
	"log"
	"os"

	mcpcli "github.com/davidl71/mcp-go-core/pkg/mcp/cli"

	"github.com/davidl71/devwisdom-go/internal/cli"
	"github.com/davidl71/devwisdom-go/internal/mcp"
)

const version = "0.1.0"

func main() {
	// Use mcp-go-core's standardized TTY detection.
	if !mcpcli.IsTTY() {
		// MCP server mode.
		server := mcp.NewWisdomServer()
		if err := server.Run(context.Background(), os.Stdin, os.Stdout); err != nil {
			log.Fatalf("Server error: %v", err)
		}
		return
	}

	// CLI mode.
	app := cli.NewApp(version)
	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
