GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=goblin

all: test build
build:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v
test:
	$(GOCMD) test ./...
clean:
	$(GOCLEAN)
	rm -rf bin/
deps:
	$(GOCMD) mod download