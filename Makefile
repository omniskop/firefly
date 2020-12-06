
PACKAGE=./cmd/firefly
BINARY_NAME=firefly

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

QTRCC=qtrcc

QTDEPLOY=qtdeploy

QT_ENV=QT_DIR="${HOME}/Qt" QT_VERSION="5.13.1" QT_API="5.13.0"
BUILD_ENV=CGO_CXXFLAGS="-g -O2 -D QT_NO_DEPRECATED_WARNINGS"
PLATFORM_WINDOWS=windows_64_shared
PLATFORM_LINUX=linux
PLATFORM_DARWIN="darwin"

all: rcc moc build

build:
	$(QT_ENV) $(BUILD_ENV) $(GOBUILD) -mod=vendor -o $(BINARY_NAME) -v $(PACKAGE)

deploy_windows: rcc_windows
	$(QTDEPLOY) -docker build $(PLATFORM_WINDOWS) $(PACKAGE)

deploy_linux: rcc_linux
	$(QTDEPLOY) -docker build $(PLATFORM_LINUX) $(PACKAGE)

deploy: rcc
	$(QT_ENV) $(QTDEPLOY) build desktop $(PACKAGE)

deploy_all: deploy deploy_windows deploy_linux

rcc:
	$(QT_ENV) qtrcc desktop $(PACKAGE)

rcc_windows:
	$(QT_ENV) qtrcc -docker $(PLATFORM_WINDOWS) $(PACKAGE)

rcc_linux:
	$(QT_ENV) qtrcc -docker $(PLATFORM_LINUX) $(PACKAGE)

moc:
	$(QT_ENV) qtmoc desktop $(PACKAGE)

setup:
	$(QT_ENV) qtsetup

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	./$(BINARY_NAME)
