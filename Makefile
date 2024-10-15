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
MAIN_PACKAGE_WEB_PATH = ./web/app

# Main binary name
BINARY_API_NAME = charmander
BINARY_WEB_NAME = charmeleon

# air tmp folder name
AIR_TMP = ./tmp

# Build api target
build:
	$(GOBUILD) -o $(BINARY_API_NAME) $(MAIN_PACKAGE_PATH)

# Build web target
build-web:
	$(GOBUILD) -o $(BINARY_WEB_NAME) $(MAIN_PACKAGE_WEB_PATH)

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
	rm -f $(BINARY_API_NAME)
	rm -f $(COVERAGE_FILE)
	rm -f $(AIR_TMP)

# Install dependencies
deps:
	$(GOGET) ./...

# Default target
all: deps fmt vet test build

# Run target (build and run)
run-api: build
	./$(BINARY_API_NAME)

# Run target (build and run)
run-web: build-web
	./$(BINARY_WEB_NAME)

# Phony targets
.PHONY: build test fmt vet clean deps all run