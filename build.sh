#!/bin/sh

export QT_DIR="${HOME}/Qt"
export QT_VERSION="5.13.1"
export QT_API="5.13.0"

package="./cmd/firefly"
binary_name="firefly"
platform_windows="windows_64_shared"
platform_linux="linux"
platform_darwin="darwin"

build()
{
    echo ">> build"
    # calling go build directly is slightly faster (one second on my machine which rougly equals 8%)
    export CGO_CXXFLAGS="-g -O2 -D QT_NO_DEPRECATED_WARNINGS"
    go build -ldflags '-s -w' -mod=vendor -o $binary_name -v $package
    # qtdeploy -fast build desktop $package
}

deploy()
{
    case $1 in
        "all") deploy_windows; deploy_linux; deploy_macos
        ;;
        "windows") deploy_windows
        ;;
        "linux") deploy_linux
        ;;
        "macos") deploy_macos
        ;;
        *) echo ">> deploy";
           qtdeploy build desktop $package
        ;;
    esac    
}

deploy_windows()
{
    echo ">> deploy for windows"
    qtdeploy -docker build $platform_windows $package
}

deploy_linux()
{
    echo ">> deploy for linux"
	qtdeploy -docker build $platform_linux $package
}

deploy_macos()
{
    echo ">> deploy for macOS"
    qtdeploy -vagrant build $platform_darwin $package
}

rcc()
{
    echo ">> resource compiler"
	qtrcc desktop $package
}

moc()
{
    echo ">> meta object compiler"
	qtmoc desktop $package
}

minimal()
{
    echo ">> minimal"
    qtminimal desktop $package
}

setup()
{
    echo ">> qt setup"
	qtsetup
}

clean()
{
    echo ">> clean"
	go clean
	rm $binary_name
}

help()
{
    echo "Firefly build script"
    echo ""
    echo "Usage:"
    echo "  build.sh               Build for the local platform."
    echo "  build.sh <command>"
    echo ""
    echo "Commands:"
    echo "  build                  Build for the local platform."
    echo "  deploy [<platform>]    Deploy Firefly for a specific platform."
    echo "  rcc                    Run the resource compiler."
    echo "  moc                    Run the meta object compiler."
    echo "  minimal                Run qtminimal."
    echo "  setup                  Setup qt bindings."
    echo "  clean                  Clean go build files."
    echo "  help                   Print this help."
    echo ""
    echo "Supported Platforms:"
    echo "  windows"
    echo "  linux"
    echo "  macos"
}

case $1 in
    "build") build
    ;;
    "deploy") deploy $2
    ;;
    "rcc") rcc
    ;;
    "moc") moc
    ;;
    "minimal") minimal
    ;;
    "setup") setup
    ;;
    "clean") clean
    ;;
    "") build
    ;;
    "help" | "h" | "--help" | "-h" | "-help" | "--h") help
    ;;
    *) echo "unknown command \"$1\"";
       echo "use help for usage"
    ;;
esac
