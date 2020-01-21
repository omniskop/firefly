
PACKAGE=./cmd/firefly
BINARY_NAME=firefly

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

QTDEPLOY=qtdeploy

ENV=CGO_CXXFLAGS="-g -O2 -D QT_NO_DEPRECATED_WARNINGS" QT_DIR="$HOME/Qt" QT_VERSION="5.13.1" QT_API="5.13.0"
WINDOWS_FLAGS=windows_64_shared
LINUX_FLAGS=linux

all: build
build:
	$(ENV) $(GOBUILD) -mod=vendor -o $(BINARY_NAME) -v $(PACKAGE)

build_win:
	$(QTDEPLOY) -docker build $(WINDOWS_FLAGS) $(PACKAGE)

build_linux:
	$(QTDEPLOY) -docker build $(LINUX_FLAGS) $(PACKAGE)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	./$(BINARY_NAME)
