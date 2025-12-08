package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/davidl71/devwisdom-go/internal/cli"
	"github.com/davidl71/devwisdom-go/internal/mcp"
)

const version = "0.1.0"

func main() {
	// Detect mode: if stdin is not a TTY, run as MCP server
	stat, _ := os.Stdin.Stat()
	isTTY := (stat.Mode() & os.ModeCharDevice) != 0

	if !isTTY {
		// MCP server mode
		server := mcp.NewWisdomServer()
		if err := server.Run(context.Background(), os.Stdin, os.Stdout); err != nil {
			log.Fatalf("Server error: %v", err)
		}
		return
	}

	// CLI mode
	app := cli.NewApp(version)
	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
