.PHONY: build build-cli build-all run test clean install install-cli fmt lint docs

# Build binary (MCP server)
build:
	go build -o devwisdom ./cmd/server

# Build CLI binary
build-cli:
	go build -o devwisdom-cli ./cmd/cli

# Build both server and CLI
build-all: build build-cli

# Run server
run: build
	./devwisdom

# Run server with watchdog (crash monitoring + file watching)
watchdog:
	./watchdog.sh --watch-files

# Run server with watchdog (restart on file changes)
watchdog-restart:
	./watchdog.sh --watch-files --restart-on-change

# Run server with watchdog (crash monitoring only)
watchdog-monitor:
	./watchdog.sh

# Run tests
test:
	go test ./...

# Run tests with coverage (text output)
test-coverage:
	go test ./... -cover -coverprofile=coverage.out
	go tool cover -func=coverage.out

# Run tests with coverage (HTML output)
test-html:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

test-wisdom:
	go test ./internal/wisdom/... -v -cover

# Clean build artifacts (both server and CLI)
clean:
	rm -f devwisdom devwisdom-cli
	go clean

# Install MCP server globally
install:
	go install ./cmd/server

# Install CLI globally
install-cli:
	go install ./cmd/cli

# Install both server and CLI globally
install-all: install install-cli

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run ./...

# Generate docs
docs:
	godoc -http=:6060
