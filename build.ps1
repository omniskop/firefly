param ([string] $command = "build", [string] $platform)

$package = "./cmd/firefly"
$binary_name = "firefly"
$platform_windows = "windows_64_shared"
$platform_linux = "linux"
$platform_darwin = "darwin"

$qtArguments = "-qt_dir=$Home/Qt", "-qt_version=5.13.1", "-qt_api=5.13.0"

function Build {
    Write-Host ">> build"
    # qtdeploy -qt_api "5.13.0" -qt_dir "$Home/Qt" -qt_version "5.13.1" -fast build desktop $package
    qtdeploy @qtArguments -fast build desktop $package
}

function Deploy {
    switch ( $platform )
    {
        "all" { DeployWindows; DeployLinux; DeployMacOS }
        "windows" { DeployWindows }
        "linux" { DeployLinux }
        "macos" { DeployMacOS }
        "" { 
            qtdeploy @qtArguments build desktop $package
        }
        default {
            Write-Host "Unsupported platform ""$platform""."
            Write-Host "Use help to get a list of supported platforms." -ErrorAction Stop
        }
    }
}

function DeployWindows {
    Write-Host ">> deploy for windows"
    qtdeploy @qtArguments -docker build $platform_windows $package
}

function DeployLinux {
    Write-Host ">> deploy for linux"
    qtdeploy @qtArguments -docker build $platform_windows $package
}

function DeployMacOS {
    Write-Host ">> deploy for macOS"
    qtdeploy @qtArguments -vagrant build $platform_darwin $package
}

function RCC {
    Write-Host ">> resource compiler"
    qtrcc @qtArguments desktop $package
}

function MOC {
    Write-Host ">> meta object compiler"
    qtmoc @qtArguments desktop $package
}

function Minimal {
    Write-Host ">> minimal"
    qtminimal @qtArguments desktop $package
}

function Setup {
    Write-Host ">> qt setup"
    qtsetup @qtArguments
}

function Clean {
    Write-Host ">> clean"
    go clean
    Remote-Item $binary_name
}

function Help {
    Write-Host "Firefly build script"
    Write-Host ""
    Write-Host "Usage:"
    Write-Host "  build.ps1              Build for the local platform."
    Write-Host "  build.ps1 <command>"
    Write-Host ""
    Write-Host "Commands:"
    Write-Host "  build                  Build for the local platform."
    Write-Host "  deploy [<platform>]    Deploy Firefly for a specific platform."
    Write-Host "  rcc                    Run the resource compiler."
    Write-Host "  moc                    Run the meta object compiler."
    Write-Host "  minimal                Run qtminimal."
    Write-Host "  setup                  Setup qt bindings."
    Write-Host "  clean                  Clean go build files."
    Write-Host "  help                   Print this help."
    Write-Host ""
    Write-Host "Supported Platforms:"
    Write-Host "  windows"
    Write-Host "  linux"
    Write-Host "  macos"
}

switch ( $command )
{
    "build" { Build }
    "deploy" { Deploy }
    "rcc" { RCC }
    "moc" { MOC }
    "minimal" { Minimal }
    "setup" { Setup }
    "clean" { Clean }
    "" { Build }
    { @("help", "h", "--help", "-h", "-help", "--h") -contains $_ } { Help }
    default {
        Write-Host "Unknown command ""$command""."
        Write-Host "Use help for usage."
    }
}
