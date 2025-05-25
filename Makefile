.PHONY: build run clean templ install dev

# Install dependencies
install:
	go mod tidy
	@echo "Installing templ CLI..."
	go install github.com/a-h/templ/cmd/templ@latest
	@echo "Checking if templ is in PATH..."
	@which templ || echo "Please add $(go env GOPATH)/bin to your PATH"

# Generate templ templates
templ:
	templ generate

# Build the application
build: templ
	go build -o bin/waqti ./main.go

# Run the application
run: templ
	go run main.go

# Development mode with auto-reload (requires air)
dev:
	@echo "Install air first: go install github.com/cosmtrek/air@latest"
	air

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf web/templates/*_templ.go

# Database setup (for future use)
db-setup:
	@echo "Database setup will be added when we integrate PostgreSQL"

# Test
test:
	go test ./...
