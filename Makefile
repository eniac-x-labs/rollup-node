# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Main package path
MAIN_PATH=./cmd

# Binary names
BINARY_NAME=rollupNode

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME) rollup-node --rpcAddress localhost:9000  --apiAddress localhost:9001

deps:
	$(GOGET) -v ./...


.PHONY: all build test clean run deps

