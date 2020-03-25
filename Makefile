
PACKAGE=./cmd/firefly
BINARY_NAME=firefly

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

QTRCC=qtrcc

QTDEPLOY=qtdeploy

ENV=CGO_CXXFLAGS="-g -O2 -D QT_NO_DEPRECATED_WARNINGS" QT_DIR="${HOME}/Qt" QT_VERSION="5.13.1" QT_API="5.13.0"
PLATFORM_WINDOWS=windows_64_shared
PLATFORM_LINUX=linux
PLATFORM_DARWIN="darwin"

build:
	$(ENV) $(GOBUILD) -mod=vendor -o $(BINARY_NAME) -v $(PACKAGE)

deploy_windows: rcc_windows
	$(QTDEPLOY) -docker build $(PLATFORM_WINDOWS) $(PACKAGE)

deploy_linux: rcc_linux
	$(QTDEPLOY) -docker build $(PLATFORM_LINUX) $(PACKAGE)

deploy: rcc
	$(ENV) $(QTDEPLOY) build desktop $(PACKAGE)

deploy_all: deploy deploy_windows deploy_linux

rcc:
	$(ENV) qtrcc desktop $(PACKAGE)

rcc_windows:
	$(ENV) qtrcc -docker $(PLATFORM_WINDOWS) $(PACKAGE)

rcc_linux:
	$(ENV) qtrcc -docker $(PLATFORM_LINUX) $(PACKAGE)

setup:
	$(ENV) qtsetup

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	./$(BINARY_NAME)
