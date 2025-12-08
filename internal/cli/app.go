package cli

import (
	"fmt"
	"os"
)

// App represents the CLI application
type App struct {
	version string
}

// NewApp creates a new CLI application instance
func NewApp(version string) *App {
	return &App{
		version: version,
	}
}

// Run executes the CLI application with the given arguments
func (a *App) Run(args []string) error {
	if len(args) == 0 {
		a.printUsage()
		return nil
	}

	command := args[0]
	commandArgs := args[1:]

	switch command {
	case "quote":
		return a.runQuote(commandArgs)
	case "consult":
		return a.runConsult(commandArgs)
	case "sources":
		return a.runSources(commandArgs)
	case "advisors":
		return a.runAdvisors(commandArgs)
	case "briefing":
		return a.runBriefing(commandArgs)
	case "version", "-v", "--version":
		fmt.Printf("devwisdom version %s\n", a.version)
		return nil
	case "help", "-h", "--help":
		a.printUsage()
		return nil
	default:
		return fmt.Errorf("unknown command: %s. Use 'devwisdom help' for usage", command)
	}
}

// printUsage prints the help message
func (a *App) printUsage() {
	fmt.Fprintf(os.Stderr, `devwisdom - Wisdom quotes and advisor consultations

USAGE:
    devwisdom <command> [options]

COMMANDS:
    quote       Get a wisdom quote
    consult     Consult an advisor
    sources     List available wisdom sources
    advisors    List available advisors
    briefing    Get daily briefing
    version     Show version
    help        Show this help message

EXAMPLES:
    devwisdom quote
    devwisdom quote --source stoic --score 75
    devwisdom consult --metric security --score 40
    devwisdom sources
    devwisdom briefing --days 7

For more information, see: https://github.com/davidl71/devwisdom-go
`)
}
