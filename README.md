# devwisdom-go

**Wisdom Module Extraction (Go Proof of Concept)**

A standalone Go MCP server providing wisdom quotes, trusted advisors, and inspirational guidance for developers. Extracted from the exarp project as a proof of concept for using compiled languages (Go) for exarp modules.

## 🎯 Project Status

**Phase 1**: ✅ Complete (Core Structure)  
**Phase 2**: ✅ Complete (Wisdom Data Porting - 16/21 local sources)  
**Phase 6**: ✅ Complete (Random Source Selector)  
**Current Phase**: Phase 3 (Advisor System)  
**Language**: Go 1.21+  
**Type**: MCP Server / Developer Tools

## 📋 Quick Start

```bash
# Clone the repository
git clone <repository-url>
cd devwisdom-go

# Build MCP server
make build

# Build CLI
make build-cli

# Run MCP server
make run

# Run CLI commands
./devwisdom-cli quote
./devwisdom-cli consult --metric security --score 75
./devwisdom-cli briefing --days 7

# Run with watchdog (crash monitoring + file watching)
make watchdog

# Test
make test
```

## 💻 CLI Usage

The `devwisdom` CLI provides easy access to wisdom quotes and advisor consultations.

### Commands

```bash
# Get a random wisdom quote
devwisdom quote

# Get quote from specific source
devwisdom quote --source stoic

# Get quote with score context
devwisdom quote --source art_of_war --score 85

# Consult an advisor for a metric
devwisdom consult --metric security --score 40

# Consult an advisor for a tool
devwisdom consult --tool project_scorecard --score 75

# Consult an advisor for a workflow stage
devwisdom consult --stage daily_checkin

# List all available wisdom sources
devwisdom sources

# List all available advisors
devwisdom advisors

# Get daily briefing
devwisdom briefing

# Get briefing for last 7 days
devwisdom briefing --days 7

# JSON output format
devwisdom quote --json
devwisdom sources --json
```

### Examples

```bash
# Quick daily wisdom
devwisdom quote

# Security advice when score is low
devwisdom consult --metric security --score 25

# Planning stage advice
devwisdom consult --stage planning --score 60

# View all Stoic quotes available
devwisdom sources --json | jq '.[] | select(.name == "stoic")'
```

## 🐚 Zsh Plugin

Install the zsh plugin for convenient shell integration.

### Installation

```bash
# Install plugin
cd zsh
./install.sh
```

Or manually:
- **oh-my-zsh**: Copy to `~/.oh-my-zsh/custom/plugins/devwisdom/`
- **Standard zsh**: Copy to `~/.zsh/plugins/devwisdom/`

Add to your `~/.zshrc`:
- **oh-my-zsh**: `plugins=(... devwisdom)`
- **Standard zsh**: `source ~/.zsh/plugins/devwisdom/devwisdom.plugin.zsh`

### Plugin Commands

```bash
# Show daily wisdom
devwisdom-daily

# Quick quote with score
devwisdom-quote stoic 75

# Consult advisor
devwisdom-consult security 40

# List sources
devwisdom-sources

# List advisors
devwisdom-advisors

# Daily briefing
devwisdom-briefing 7
```

### Auto-Daily Wisdom

Enable automatic daily wisdom on shell startup:

```bash
# Add to ~/.zshrc
export DEVWISDOM_AUTO_DAILY=true
```

## 🔄 Watchdog Script

The project includes a watchdog script that monitors the server for crashes and can automatically reload on file changes.

### Basic Usage

```bash
# Monitor for crashes only
./watchdog.sh

# Watch files and reload on changes
./watchdog.sh --watch-files

# Watch files and restart on changes
./watchdog.sh --watch-files --restart-on-change

# Log to file
./watchdog.sh --watch-files --log-file watchdog.log
```

### Features

- **Crash Detection**: Automatically restarts server if it crashes
- **File Watching**: Monitors `sources.json`, config files, and Go source files
- **Hot Reload**: Reloads sources without full restart (if supported)
- **Auto Rebuild**: Rebuilds binary when Go source files change
- **Restart Limits**: Configurable maximum restart attempts
- **Graceful Shutdown**: Handles SIGTERM for clean shutdowns

### Watchdog Options

- `--watch-files`: Enable file watching (uses `fswatch` on macOS, falls back to polling)
- `--restart-on-change`: Restart server on file changes (default: reload)
- `--log-file PATH`: Log output to file
- `--pid-file PATH`: Custom PID file location
- `--max-restarts N`: Maximum restart attempts (default: 10)
- `--restart-delay SEC`: Delay between restarts (default: 2)

### Makefile Targets

```bash
make watchdog          # Run with file watching and reload
make watchdog-restart  # Run with file watching and restart on change
make watchdog-monitor  # Run with crash monitoring only
```

## 🏗️ Project Structure

```
devwisdom-go/
├── cmd/server/          # MCP server entry point
├── internal/
│   ├── wisdom/         # Wisdom engine (quotes, sources, advisors)
│   ├── mcp/            # MCP protocol handler
│   └── config/         # Configuration management
├── docs/               # Documentation
├── Makefile           # Build commands
└── go.mod             # Go dependencies
```

## 📊 Planning & Status

**Todo2 Tasks**: 37 tasks across 9 phases (tracked in agentic-tools MCP)  
**Project ID**: `039bb05a-6f78-492b-88b5-28fdfa3ebce7`

See `PROJECT_GOALS.md` for detailed phase breakdown and `PRD.md` for full requirements.

## 🚀 Phases

1. ✅ **Phase 1**: Core Structure (Complete)
2. ✅ **Phase 2**: Wisdom Data Porting (16/21 local sources complete)
3. ⏳ **Phase 3**: Advisor System
4. ⏳ **Phase 4**: MCP Protocol Implementation
5. ⏳ **Phase 5**: Consultation Logging
6. ✅ **Phase 6**: Daily Random Source Selection (Complete)
7. ⏳ **Phase 7**: Optional Features (Sefaria, TTS)
8. ⏳ **Phase 8**: Testing
9. ⏳ **Phase 9**: Documentation
10. ⏳ **Phase 10**: Polish & Deployment

## 📚 Documentation

- **PROJECT_GOALS.md** - Strategic phases and goals
- **PRD.md** - Product Requirements Document (129 user stories)
- **TODO.md** - Task breakdown by phase
- **EXARP_PLANNING_COMPLETE.md** - Planning analysis summary
- **docs/WATCHDOG.md** - Watchdog script documentation
- **docs/CURSOR_EXTENSION.md** - Cursor extension architecture (future enhancement)

## 🔗 Related

- **Source**: Python wisdom module in `exarp` project
- **MCP Spec**: https://modelcontextprotocol.io/
- **Go Docs**: https://go.dev/doc/effective_go

## 📝 License

[Add your license here]

## 👤 Author

Extracted from exarp project as compiled language PoC.
