# Watchdog Script Documentation

## Overview

The `watchdog.sh` script provides automatic crash recovery and file change monitoring for the devwisdom-go MCP server. It ensures the server stays running and can automatically reload or restart when configuration or source files change.

## Features

### Crash Monitoring
- Automatically detects when the server process crashes
- Restarts the server with configurable retry limits
- Tracks restart attempts to prevent infinite restart loops
- Graceful shutdown handling (SIGTERM before SIGKILL)

### File Watching
- Monitors `sources.json` and configuration files
- Watches Go source files for changes
- Automatically rebuilds binary when Go files change
- Supports reload or restart on file changes

### Platform Support
- **macOS**: Uses `fswatch` for efficient file watching (install with `brew install fswatch`)
- **Fallback**: Polling-based watching if `fswatch` is not available
- Works on Linux with `inotifywait` (modify script for Linux)

## Usage

### Basic Crash Monitoring

```bash
# Monitor for crashes only (no file watching)
./watchdog.sh
```

### File Watching with Reload

```bash
# Watch files and reload sources on change
./watchdog.sh --watch-files
```

### File Watching with Restart

```bash
# Watch files and restart server on change
./watchdog.sh --watch-files --restart-on-change
```

### Advanced Options

```bash
# Custom log file
./watchdog.sh --watch-files --log-file /var/log/devwisdom-watchdog.log

# Custom PID file location
./watchdog.sh --pid-file /tmp/devwisdom.pid

# Limit restart attempts
./watchdog.sh --max-restarts 5

# Custom restart delay
./watchdog.sh --restart-delay 5
```

### Makefile Targets

```bash
make watchdog          # File watching + reload on change
make watchdog-restart  # File watching + restart on change
make watchdog-monitor  # Crash monitoring only
```

## Configuration

### Watched Files

The script watches the following patterns:
- `sources.json` (root)
- `.wisdom/sources.json` (project-specific)
- `wisdom/sources.json` (alternative location)
- `*.go` (all Go source files)
- `go.mod` and `go.sum` (dependency files)

### Default Settings

- **Max Restarts**: 10 attempts
- **Restart Delay**: 2 seconds
- **PID File**: `.devwisdom.pid` (project root)
- **Binary**: `devwisdom` (project root)

## How It Works

### Process Monitoring

1. Checks if server process is running every 2 seconds
2. If process dies, attempts restart
3. Tracks restart count to prevent infinite loops
4. Rebuilds binary if Go source files changed

### File Watching

**With fswatch (macOS):**
- Uses `fswatch` for efficient event-driven file watching
- Triggers reload/restart immediately on file change

**Polling Fallback:**
- Checks file modification times every 2 seconds
- Less efficient but works without external dependencies

### Reload vs Restart

**Reload** (default):
- Sends USR1 signal to server (if supported)
- Falls back to restart if signal not supported
- Faster, preserves server state

**Restart**:
- Stops and starts server process
- Ensures clean state
- Rebuilds binary if needed

## Signal Handling

The watchdog handles:
- **SIGINT** (Ctrl+C): Graceful shutdown
- **SIGTERM**: Graceful shutdown
- **SIGUSR1**: Reload signal (sent to server if supported)

## Logging

### Console Output

Colored output for different log levels:
- 🟢 **INFO**: Normal operations
- 🟡 **WARN**: Warnings (e.g., server already running)
- 🔴 **ERROR**: Errors (crashes, failures)
- 🔵 **DEBUG**: Debug information

### File Logging

When `--log-file` is specified:
- All output is written to the log file
- Also displayed on console (tee behavior)
- Useful for production deployments

## Troubleshooting

### Server Won't Start

1. Check if binary exists: `ls -la devwisdom`
2. Build manually: `make build`
3. Check permissions: `chmod +x devwisdom`
4. Check logs for errors

### File Changes Not Detected

1. Install fswatch: `brew install fswatch`
2. Check if files are in watched paths
3. Verify file permissions
4. Try polling mode (remove fswatch)

### Too Many Restarts

1. Check server logs for crash reasons
2. Increase `--restart-delay`
3. Lower `--max-restarts` to fail faster
4. Fix underlying server issues

### PID File Issues

1. Remove stale PID file: `rm .devwisdom.pid`
2. Check if process is actually running: `ps aux | grep devwisdom`
3. Use custom PID file location: `--pid-file /tmp/custom.pid`

## Integration with Cursor MCP

When using with Cursor's MCP integration:

1. Start watchdog in background: `./watchdog.sh --watch-files &`
2. Configure Cursor to use the server
3. Watchdog will handle crashes and reloads automatically
4. Changes to `sources.json` will be picked up automatically

## Example Workflow

```bash
# Terminal 1: Start watchdog
cd /path/to/devwisdom-go
make watchdog

# Terminal 2: Edit sources.json
vim sources.json
# Add new source...

# Watchdog automatically detects change and reloads
# Server picks up new source without manual restart
```

## Future Enhancements

Potential improvements:
- [ ] Signal-based reload support in server (USR1 handler)
- [ ] Webhook notifications on crashes
- [ ] Metrics collection (uptime, restart count)
- [ ] Health check endpoint
- [ ] Docker container support
- [ ] Systemd service file generation

