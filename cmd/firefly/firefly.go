package main

import (
	"os"
	"runtime"

	"github.com/therecipe/qt/core"

	"github.com/sirupsen/logrus"
	"github.com/therecipe/qt/widgets"
)

var app *widgets.QApplication

func init() {
	// not sure if this is really needed
	runtime.LockOSThread()
}

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.Info("FireFly starting...")
	wd, _ := os.Getwd()
	exe, _ := os.Executable()
	logrus.WithFields(logrus.Fields{"currentPath": core.QDir_CurrentPath(), "wd": wd, "exe": exe}).Info("Directories")

	app = createApplication()

	logrus.Info("Application created")

	OpenLaunchwindow()

	logrus.Info("Starting Application")

	app.Exec()
}
