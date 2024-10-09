# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOTOOLCOVER = $(GOCMD) tool cover
GOFMT = $(GOCMD) fmt
GOVET = $(GOCMD) vet

COVERAGE_FILE = coverage.out

# Main package path
MAIN_PACKAGE_PATH = ./cmd/app

# Main binary name
BINARY_NAME = charmander

# air tmp folder name
AIR_TMP = ./tmp

# Build target
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PACKAGE_PATH)

# Test target
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -cover ./...

# Generate coverage profile
coverage-profile:
	$(GOTEST) -coverprofile=$(COVERAGE_FILE) ./...

# View coverage in browser
coverage-html: coverage-profile
	$(GOTOOLCOVER) -html=$(COVERAGE_FILE)

# Run tests and generate a coverage report
test-with-coverage: coverage-profile
	$(GOTOOLCOVER) -func=$(COVERAGE_FILE)

# Format source code
fmt:
	$(GOFMT) ./...

# Vet source code
vet:
	$(GOVET) ./...

# Clean target
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(COVERAGE_FILE)
	rm -f $(AIR_TMP)

# Install dependencies
deps:
	$(GOGET) ./...

# Default target
all: deps fmt vet test build

# Run target (build and run)
run: build
	./$(BINARY_NAME)

# Phony targets
.PHONY: build test fmt vet clean deps all run