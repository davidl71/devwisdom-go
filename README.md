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

The `devwisdom` CLI provides easy access to wisdom quotes and advisor consultations. The CLI can run in two modes:
- **CLI Mode**: When run from a terminal (stdin is a TTY)
- **MCP Server Mode**: When run via stdio (for MCP integration)

### Building the CLI

```bash
# Build CLI binary
make build-cli

# Build both server and CLI
make build-all

# Install CLI globally
make install-cli
```

### Commands

```bash
# Get a random wisdom quote (date-seeded, consistent per day)
devwisdom quote

# Get quote from specific source
devwisdom quote --source stoic

# Get quote with score context (affects aeon level selection)
devwisdom quote --source art_of_war --score 85

# Get quote in JSON format
devwisdom quote --json

# Get quote text only (quiet mode)
devwisdom quote --quiet

# Consult an advisor for a metric
devwisdom consult --metric security --score 40

# Consult an advisor for a tool
devwisdom consult --tool project_scorecard --score 75

# Consult an advisor for a workflow stage
devwisdom consult --stage daily_checkin --score 60

# List all available wisdom sources
devwisdom sources

# List sources in JSON format
devwisdom sources --json

# List all available advisors
devwisdom advisors

# List advisors in JSON format
devwisdom advisors --json

# Get daily briefing (today's wisdom)
devwisdom briefing

# Get briefing for last 7 days
devwisdom briefing --days 7

# Get briefing in JSON format
devwisdom briefing --json

# Show version
devwisdom version

# Show help
devwisdom help
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

# Get quote for high-performing project
devwisdom quote --source pistis_sophia --score 90

# Get quiet quote for scripts
QUOTE=$(devwisdom quote --quiet)
echo "Today: $QUOTE"
```

### Command Options

**`quote` command:**
- `--source SOURCE`: Wisdom source name (e.g., `stoic`, `pistis_sophia`, `random`)
- `--score SCORE`: Project health score (0-100), affects aeon level selection
- `--json`: Output in JSON format
- `--quiet`: Output only the quote text

**`consult` command:**
- `--metric METRIC`: Metric name (e.g., `security`, `testing`, `documentation`)
- `--tool TOOL`: Tool name (e.g., `project_scorecard`)
- `--stage STAGE`: Stage name (e.g., `daily_checkin`, `sprint_planning`)
- `--score SCORE`: Project health score (0-100), required for metric/tool consultations
- `--json`: Output in JSON format
- `--quiet`: Output only the quote text

**`briefing` command:**
- `--days DAYS`: Number of days to include (default: 1)
- `--json`: Output in JSON format

### Use Cases

**Daily Standup:**
```bash
# Get today's wisdom for team standup
devwisdom quote
```

**Project Health Monitoring:**
```bash
# Get advice based on current project score
SCORE=$(calculate-project-score)
devwisdom consult --metric security --score $SCORE
```

**CI/CD Integration:**
```bash
# Add wisdom to build notifications
QUOTE=$(devwisdom quote --quiet)
echo "Build wisdom: $QUOTE" | send-notification
```

**Documentation Generation:**
```bash
# Generate markdown with quotes
devwisdom sources --json | jq -r '.[] | "## \(.name)\n\n\(.description)\n"'
```

📚 **For detailed usage examples, see:**
- [CLI Usage Examples](examples/cli_usage.md) - Comprehensive CLI command examples with output samples
- [MCP Integration Guide](examples/mcp_integration.md) - MCP server integration examples
- [Programmatic API Examples](examples/programmatic_usage.go) - Go library usage examples
- [Examples Directory](examples/) - All example files and guides

### Quick Reference

| Command | Description | Example |
|---------|-------------|---------|
| `quote` | Get wisdom quote | `devwisdom quote --source stoic` |
| `consult` | Consult advisor | `devwisdom consult --metric security --score 40` |
| `sources` | List sources | `devwisdom sources` |
| `advisors` | List advisors | `devwisdom advisors` |
| `briefing` | Daily briefing | `devwisdom briefing --days 7` |
| `version` | Show version | `devwisdom version` |
| `help` | Show help | `devwisdom help` |

## 🐚 Zsh Plugin

Install the zsh plugin for convenient shell integration with tab completion and helper functions.

### Installation

#### Automatic Installation (Recommended)

```bash
# From project root
cd zsh
./install.sh
```

The installation script will:
- Detect your zsh framework (oh-my-zsh or standard zsh)
- Copy plugin files to the appropriate location
- Configure your `~/.zshrc` automatically
- Set up tab completion

#### Manual Installation

**For oh-my-zsh users:**

```bash
# Create plugin directory
mkdir -p ~/.oh-my-zsh/custom/plugins/devwisdom

# Copy plugin files
cp zsh/devwisdom.plugin.zsh ~/.oh-my-zsh/custom/plugins/devwisdom/
cp zsh/_devwisdom ~/.oh-my-zsh/custom/plugins/devwisdom/

# Make completion executable
chmod +x ~/.oh-my-zsh/custom/plugins/devwisdom/_devwisdom

# Add to ~/.zshrc
echo 'plugins=(... devwisdom)' >> ~/.zshrc
```

**For standard zsh users:**

```bash
# Create plugin directory
mkdir -p ~/.zsh/plugins/devwisdom

# Copy plugin files
cp zsh/devwisdom.plugin.zsh ~/.zsh/plugins/devwisdom/
cp zsh/_devwisdom ~/.zsh/plugins/devwisdom/

# Make completion executable
chmod +x ~/.zsh/plugins/devwisdom/_devwisdom

# Add to ~/.zshrc
cat >> ~/.zshrc << 'EOF'
# devwisdom plugin
source ~/.zsh/plugins/devwisdom/devwisdom.plugin.zsh
fpath=(~/.zsh/plugins/devwisdom $fpath)
autoload -Uz compinit && compinit
EOF
```

**After installation:**
```bash
# Reload your shell configuration
source ~/.zshrc
```

### Plugin Commands

The zsh plugin provides convenient wrapper functions:

```bash
# Show daily wisdom (random source, date-seeded)
devwisdom-daily

# Show daily wisdom from specific source
devwisdom-daily stoic

# Quick quote with optional source and score
devwisdom-quote              # Random quote
devwisdom-quote stoic        # Quote from stoic source
devwisdom-quote stoic 75     # Quote from stoic with score 75

# Consult advisor (metric, tool, or stage)
devwisdom-consult security 40                    # Metric advisor
devwisdom-consult project_scorecard 75          # Tool advisor
devwisdom-consult daily_checkin 60               # Stage advisor

# List available sources
devwisdom-sources

# List available advisors
devwisdom-advisors

# Get daily briefing (default: 1 day)
devwisdom-briefing
devwisdom-briefing 7         # Last 7 days
```

### Tab Completion

The plugin includes full tab completion support:

```bash
# Tab completion for commands
devwisdom <TAB>              # Shows: quote, consult, sources, advisors, briefing

# Tab completion for options
devwisdom quote --<TAB>       # Shows: --source, --score, --json, --quiet
devwisdom consult --<TAB>    # Shows: --metric, --tool, --stage, --score, --json, --quiet
```

### Auto-Daily Wisdom

Enable automatic daily wisdom on shell startup:

```bash
# Add to ~/.zshrc
export DEVWISDOM_AUTO_DAILY=true
```

When enabled, you'll see a daily wisdom quote each time you open a new terminal session.

### Use Cases

**Daily Motivation:**
```bash
# Add to ~/.zshrc for daily inspiration
devwisdom-daily
```

**Project Health Check:**
```bash
# Quick consultation when starting work
devwisdom-consult security $(get-project-score)
```

**Script Integration:**
```bash
#!/bin/zsh
# Get wisdom for automation scripts
QUOTE=$(devwisdom quote --json)
echo "Today's wisdom: $(echo $QUOTE | jq -r '.quote')"
```

**Custom Aliases:**
```bash
# Add to ~/.zshrc
alias wisdom='devwisdom quote'
alias advice='devwisdom consult'
alias daily='devwisdom-daily'
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
7. ⏳ **Phase 7**: Optional Features (Sefaria API)
8. ⏳ **Phase 8**: Testing
9. ⏳ **Phase 9**: Documentation
10. ⏳ **Phase 10**: Polish & Deployment

## 📚 Documentation

- **PROJECT_GOALS.md** - Strategic phases and goals
- **PRD.md** - Product Requirements Document (129 user stories)
- **TODO.md** - Task breakdown by phase
- **EXARP_PLANNING_COMPLETE.md** - Planning analysis summary
- **docs/WATCHDOG.md** - Watchdog script documentation
- **docs/CURSOR_EXTENSION.md** - Cursor extension architecture (⚠️ **Future goal - not currently implemented**)

## 🔗 Related

- **Source**: Python wisdom module in `exarp` project
- **MCP Spec**: https://modelcontextprotocol.io/
- **Go Docs**: https://go.dev/doc/effective_go

## 📝 License

[Add your license here]

## 👤 Author

Extracted from exarp project as compiled language PoC.
