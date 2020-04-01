package main

import (
	"flag"
	"os"
	"runtime"

	"github.com/omniskop/firefly/cmd/firefly/editor"
	"github.com/omniskop/firefly/pkg/storage"

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
	flag.Parse()
	wd, _ := os.Getwd()
	exe, _ := os.Executable()
	logrus.WithFields(logrus.Fields{"currentPath": core.QDir_CurrentPath(), "wd": wd, "exe": exe}).Info("Directories")

	app = createApplication()

	logrus.Info("Application created")

	if fileName := flag.Arg(0); fileName != "" {
		project, err := storage.LoadFile(fileName)
		if err != nil {
			logrus.Error(err)
			return
		}
		edit := editor.New(project, ApplicationCallbacks)
		edit.SaveLocation = fileName
	} else {
		OpenLaunchwindow()
	}

	logrus.Info("Starting Application")

	app.Exec()
}
