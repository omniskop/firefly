package main

import (
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)


func createApplication() *widgets.QApplication {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetAttribute(core.Qt__AA_UseHighDpiPixmaps, true)
	return app
}
