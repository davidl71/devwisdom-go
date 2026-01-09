.PHONY: build build-cli build-all run test clean install install-cli fmt lint lint-fix docs bench bench-cpu bench-mem bench-profile pprof-cpu pprof-mem pprof-web-cpu pprof-web-mem build-windows build-linux build-darwin build-all-platforms build-release clean-dist

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
clean: clean-dist
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

# Lint code and auto-fix issues
lint-fix:
	golangci-lint run --fix ./...

# Generate docs
docs:
	godoc -http=:6060

# Run benchmarks
bench:
	go test -bench=. -benchmem -benchtime=3s ./internal/wisdom/...

# Run benchmarks with CPU profiling
bench-cpu:
	go test -bench=. -benchmem -benchtime=3s -cpuprofile=cpu.prof ./internal/wisdom/...
	@echo "CPU profile saved to cpu.prof"
	@echo "Analyze with: go tool pprof cpu.prof"

# Run benchmarks with memory profiling
bench-mem:
	go test -bench=. -benchmem -benchtime=3s -memprofile=mem.prof ./internal/wisdom/...
	@echo "Memory profile saved to mem.prof"
	@echo "Analyze with: go tool pprof mem.prof"

# Run benchmarks with both CPU and memory profiling
bench-profile: bench-cpu bench-mem

# Analyze CPU profile (interactive)
pprof-cpu:
	go tool pprof cpu.prof

# Analyze memory profile (interactive)
pprof-mem:
	go tool pprof mem.prof

# Generate profile reports (web interface)
pprof-web-cpu:
	go tool pprof -http=:8080 cpu.prof

# Generate profile reports (web interface)
pprof-web-mem:
	go tool pprof -http=:8080 mem.prof

# Cross-compilation variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
DIST_DIR = dist
BUILD_FLAGS = -ldflags="-s -w" # Strip debug symbols and reduce binary size

# Build Windows binaries (amd64)
build-windows:
	@echo "Building Windows binaries..."
	@mkdir -p $(DIST_DIR)/windows-amd64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/windows-amd64/devwisdom.exe ./cmd/server
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/windows-amd64/devwisdom-cli.exe ./cmd/cli
	@echo "Windows binaries built in $(DIST_DIR)/windows-amd64/"

# Build Linux binaries (amd64, arm64)
build-linux:
	@echo "Building Linux binaries..."
	@mkdir -p $(DIST_DIR)/linux-amd64 $(DIST_DIR)/linux-arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/linux-amd64/devwisdom ./cmd/server
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/linux-amd64/devwisdom-cli ./cmd/cli
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/linux-arm64/devwisdom ./cmd/server
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/linux-arm64/devwisdom-cli ./cmd/cli
	@echo "Linux binaries built in $(DIST_DIR)/linux-*/"

# Build macOS binaries (amd64, arm64)
build-darwin:
	@echo "Building macOS binaries..."
	@mkdir -p $(DIST_DIR)/darwin-amd64 $(DIST_DIR)/darwin-arm64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/darwin-amd64/devwisdom ./cmd/server
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/darwin-amd64/devwisdom-cli ./cmd/cli
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/darwin-arm64/devwisdom ./cmd/server
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/darwin-arm64/devwisdom-cli ./cmd/cli
	@echo "macOS binaries built in $(DIST_DIR)/darwin-*/"

# Build all platforms
build-all-platforms: build-windows build-linux build-darwin
	@echo "All platform binaries built in $(DIST_DIR)/"

# Build and package release archives
build-release: clean-dist build-all-platforms
	@echo "Creating release archives..."
	@mkdir -p $(DIST_DIR)/release
	@# Windows release
	cd $(DIST_DIR)/windows-amd64 && zip -q ../release/devwisdom-$(VERSION)-windows-amd64.zip devwisdom.exe devwisdom-cli.exe
	@# Linux releases
	cd $(DIST_DIR)/linux-amd64 && tar czf ../release/devwisdom-$(VERSION)-linux-amd64.tar.gz devwisdom devwisdom-cli
	cd $(DIST_DIR)/linux-arm64 && tar czf ../release/devwisdom-$(VERSION)-linux-arm64.tar.gz devwisdom devwisdom-cli
	@# macOS releases
	cd $(DIST_DIR)/darwin-amd64 && tar czf ../release/devwisdom-$(VERSION)-darwin-amd64.tar.gz devwisdom devwisdom-cli
	cd $(DIST_DIR)/darwin-arm64 && tar czf ../release/devwisdom-$(VERSION)-darwin-arm64.tar.gz devwisdom devwisdom-cli
	@echo "Release archives created in $(DIST_DIR)/release/"
	@ls -lh $(DIST_DIR)/release/

# Clean distribution directory
clean-dist:
	rm -rf $(DIST_DIR)
