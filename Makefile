GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGENERATE=$(GOCMD) generate
BINARY_NAME=goblin

all: generate test build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
test:
	$(GOCMD) test ./...
clean:
	$(GOCLEAN)
	rm -rf bin/
deps:
	$(GOCMD) mod download
generate:
	$(GOGENERATE) ./...