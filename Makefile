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
MAIN_PACKAGE_TEMPL_PATH = ./templ/app

# Main binary name
BINARY_API_NAME = charmander
# BINARY_WEB_NAME = charmeleon
BINARY_TEMPL_NAME = charizard

# air tmp folder name
AIR_TMP = ./tmp

# Build api target
build:
	$(GOBUILD) -o $(BINARY_API_NAME) $(MAIN_PACKAGE_PATH)

# Build web target used with go template (not anymore)
# build-web:
# 	$(GOBUILD) -o $(BINARY_WEB_NAME) $(MAIN_PACKAGE_WEB_PATH)

build-templ-file:
	templ generate ./templ

# Build templ target
build-templ: build-templ-file
	$(GOBUILD) -o $(BINARY_TEMPL_NAME) $(MAIN_PACKAGE_TEMPL_PATH)

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

# Run api target (build and run)
run-api: build
	./$(BINARY_API_NAME)

# Run web target (build and run)
# run-web: build-web
# 	./$(BINARY_WEB_NAME)

# Run templ target (build and run)
run-templ: build-templ
	./$(BINARY_TEMPL_NAME)

# Run templ target (build and run)
run-templ-raw:
	$(GOBUILD) -o $(BINARY_TEMPL_NAME) $(MAIN_PACKAGE_TEMPL_PATH)
	./$(BINARY_TEMPL_NAME)

# Phony targets
.PHONY: build test fmt vet clean deps all run