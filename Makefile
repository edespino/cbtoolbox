# Project settings
BUILD_DIR = build
EXECUTABLE := $(BUILD_DIR)/cbtoolbox
SOURCES := $(shell go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' ./... | tr '\n' ' ')
EMBEDDED_FILES := cmd/coreinfo/resources/gdb_commands_basic.txt cmd/coreinfo/resources/gdb_commands_detailed.txt


# Go settings
GO = go
GO_TEST_FLAGS = -v
GO_LINT_TOOL = golangci-lint

.PHONY: build

# Default target
all: build

# Build the project
build: $(EXECUTABLE)

$(EXECUTABLE): $(SOURCES) $(EMBEDDED_FILES)
	mkdir -p $(BUILD_DIR)
	go build -o $(EXECUTABLE) main.go

# Run all tests
test:
	$(GO) test $(GO_TEST_FLAGS) ./...

# Run tests with coverage
test-cover:
	$(GO) test $(GO_TEST_FLAGS) -cover ./...

# Lint the code
lint:
	@command -v $(GO_LINT_TOOL) >/dev/null 2>&1 || { \
		echo "Error: $(GO_LINT_TOOL) is not installed. Install it using 'go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest'."; \
		exit 1; \
	}
	$(GO_LINT_TOOL) run

.PHONY: clean
# Clean up build artifacts
clean:
	rm -rf $(BUILD_DIR)

# Run the application
run:
	$(BUILD_DIR)/$(BINARY_NAME)

# Help
help:
	@echo "Available targets:"
	@echo "  build        Build the binary"
	@echo "  test         Run all tests"
	@echo "  test-cover   Run tests with coverage"
	@echo "  lint         Run linters"
	@echo "  clean        Clean build artifacts"
	@echo "  run          Run the application"
	@echo "  help         Show this help message"
