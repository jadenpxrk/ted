.PHONY: build clean install cross-compile release run dev help

# Version and build flags
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-s -w -X ted/cmd.Version=$(VERSION)"

# Build the binary for the current platform
build:
	go build $(LDFLAGS) -o ted

# Build for multiple platforms
cross-compile:
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/ted-darwin-arm64
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/ted-darwin-amd64
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/ted-linux-amd64
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/ted-linux-arm64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/ted-windows-amd64.exe

# Install to /usr/local/bin
install: build
	sudo cp ted /usr/local/bin/

# Clean build artifacts
clean:
	rm -f ted
	rm -rf dist/

# Run the built binary
run: build
	./ted

# Update dependencies
dev:
	go mod tidy

# Create release builds
release: clean
	mkdir -p dist
	$(MAKE) cross-compile
	cd dist && shasum -a 256 * > checksums.txt

# Show available targets
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary for current platform"
	@echo "  cross-compile- Build for all platforms"
	@echo "  install      - Install to /usr/local/bin"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Build and run the binary"
	@echo "  dev          - Update dependencies"
	@echo "  release      - Create release builds with checksums"
	@echo "  help         - Show this help"
