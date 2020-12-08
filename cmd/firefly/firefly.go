package main

import (
	"encoding/json"
	"flag"
	"os"
	"runtime"

	"github.com/omniskop/firefly/pkg/scanner"

	"github.com/omniskop/firefly/cmd/firefly/settings"

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
	logrus.AddHook(&logConsoleReceiver)
	logrus.SetLevel(logrus.TraceLevel)
	logrus.Info("FireFly starting...")
	flag.Parse()
	wd, _ := os.Getwd()
	exe, _ := os.Executable()
	logrus.WithFields(logrus.Fields{"currentPath": core.QDir_CurrentPath(), "wd": wd, "exe": exe}).Info("Directories")

	versionCheck()

	app = createApplication()

	logrus.Info("Application created")

	if fileName := flag.Arg(0); fileName != "" {
		project, err := storage.LoadFile(fileName)
		if err != nil {
			logrus.Error(err)
			return
		}
		addRecentFile(newRecentFileDetailed(project.Audio.Title, project.Audio.Author, fileName))
		edit := editor.New(project, ApplicationCallbacks)
		edit.SaveLocation = fileName
	} else {
		OpenLaunchWindow()
	}

	logrus.Info("Starting Application")

	app.Exec()
}

func versionCheck() {
	version := settings.GetString("version")

	switch version {
	case "": // new installation
		restoreDefaultSettings()
		logrus.Info("new installation: settings have been reset")
	case "0.1.0":
		mapping, _ := json.Marshal(scanner.NewLinearMapping(settings.GetInt("ledCount")))
		settings.Set("liveLedStrip/mapping", string(mapping))
		settings.Remove("ledCount")
		fallthrough
	case "0.1.1":
	}

	settings.Set("version", "0.1.1")
}
