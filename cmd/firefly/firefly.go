package main

import (
	"runtime"
	"github.com/sirupsen/logrus"
)

var app *widgets.QApplication

func init() {
	// not shure if this is realy needed
	runtime.LockOSThread()
}

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.Info("FireFly starting...")

	app = createApplication()

	logrus.Info("Application created")

	logrus.Info("Starting Application")

	app.Exec()
}
