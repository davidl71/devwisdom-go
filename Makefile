.PHONY: build run test clean install

# Build binary (MCP server)
build:
	go build -o devwisdom ./cmd/server

# Build CLI binary
build-cli:
	go build -o devwisdom-cli ./cmd/cli

# Run server
run: build
	./devwisdom

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html -v

test-coverage:
	go test ./... -cover -coverprofile=coverage.out
	go tool cover -func=coverage.out

test-html:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

test-wisdom:
	go test ./internal/wisdom/... -v -cover

# Clean build artifacts
clean:
	rm -f devwisdom
	go clean

# Install globally
install:
	go install ./cmd/server

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run ./...

# Generate docs
docs:
	godoc -http=:6060
