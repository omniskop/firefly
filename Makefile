
PACKAGE=./cmd/firefly
BINARY_NAME=firefly

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: build
build:
	$(GOBUILD) -mod=vendor -o $(BINARY_NAME) -v $(PACKAGE)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	./$(BINARY_NAME)
