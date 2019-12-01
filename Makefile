
PACKAGE=./cmd/firefly
BINARY_NAME=firefly

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

ENV=CGO_CXXFLAGS="-g -O2 -D QT_NO_DEPRECATED_WARNINGS"

all: build
build:
	$(ENV) $(GOBUILD) -mod=vendor -o $(BINARY_NAME) -v $(PACKAGE)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	./$(BINARY_NAME)
