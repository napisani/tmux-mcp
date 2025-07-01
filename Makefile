
# Variables
APP_NAME=tmux-mcp
BINARY_NAME=$(APP_NAME)
VERSION=0.1.0
BUILD_DIR=bin
SRC_DIR=cmd/$(APP_NAME)

# Build the binary
.PHONY: build
build:
	@echo "Building the application..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)/main.go

# Run the application
.PHONY: run
run: build
	@echo "Running the application..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)

# Create a distribution archive
.PHONY: dist
dist: build
	@echo "Creating distribution archive..."
	mkdir -p release
	tar -czf release/$(BINARY_NAME)-$(VERSION).tar.gz -C $(BUILD_DIR) $(BINARY_NAME)

# Watch for changes and rebuild
.PHONY: watch
watch:
	@echo "Watching for changes..."
	@nodemon  --exec "make build ; exit 0" --watch $(SRC_DIR) --ext go  --signal SIGINT  

.PHONY: inspect
inspect:
	@echo "Inspecting the application..."
	@npx   npx @modelcontextprotocol/inspector  ./bin/$(BINARY_NAME) 
# Help command
.PHONY: help
help:
	@echo "Makefile commands:"
	@echo "  make build   - Build the application"
	@echo "  make run     - Build and run the application"
	@echo "  make clean   - Clean up build artifacts"
	@echo "  make dist    - Create a distribution archive"
	@echo "  make watch   - Watch for changes and rebuild"
	@echo "  make inspect - Inspect the application with Model Context Protocol"
	@echo "  make help    - Show this help message"



