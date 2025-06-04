.PHONY: build clean test install cross-compile release

# Build the binary for the current platform
build:
	go build -o ted

# Build for all platforms
cross-compile:
	GOOS=darwin GOARCH=arm64 go build -o dist/ted-darwin-arm64
	GOOS=darwin GOARCH=amd64 go build -o dist/ted-darwin-amd64
	GOOS=linux GOARCH=amd64 go build -o dist/ted-linux-amd64
	GOOS=windows GOARCH=amd64 go build -o dist/ted-windows-amd64.exe

# Install to /usr/local/bin (requires sudo)
install: build
	sudo mv ted /usr/local/bin/

# Clean build artifacts
clean:
	rm -f ted
	rm -rf dist/

# Run tests
test:
	go test -v ./...

# Run the built binary
run: build
	./ted

# Development setup
dev-setup:
	go mod tidy
	go get -u all

# Create release directory
release: clean
	mkdir -p dist
	$(MAKE) cross-compile
	cd dist && shasum -a 256 * > checksums.txt

# Help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary for current platform"
	@echo "  cross-compile- Build for all platforms"  
	@echo "  install      - Install to /usr/local/bin"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  run          - Build and run"
	@echo "  dev-setup    - Set up development environment"
	@echo "  release      - Create release builds with checksums"
	@echo "  help         - Show this help"
