.PHONY: install-air air dev run build clean

# Install air for hot reload
install-air:
	go install github.com/cosmtrek/air@latest

# Run with hot reload
air:
	air

# Development mode with hot reload
dev: 
	@which air > /dev/null || make install-air
	air

# Run without hot reload
run:
	go run cmd/server/main.go

# Build binary
build:
	go build -o bin/server cmd/server/main.go

# Clean temporary files
clean:
	rm -rf tmp/
	rm -rf bin/