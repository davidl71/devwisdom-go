.PHONY: build run test clean install

# Build binary
build:
	go build -o devwisdom ./cmd/server

# Run server
run: build
	./devwisdom

# Run tests
test:
	go test ./...

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
