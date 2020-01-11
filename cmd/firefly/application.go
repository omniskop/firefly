package main

import (
	"os"

	"github.com/omniskop/firefly/cmd/firefly/editor"
	"github.com/omniskop/firefly/pkg/storage"
	"github.com/sirupsen/logrus"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func createApplication() *widgets.QApplication {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetAttribute(core.Qt__AA_UseHighDpiPixmaps, true)
	return app
}

var ApplicationCallbacks map[string]func()

func init() {
	ApplicationCallbacks = map[string]func(){
		"open": func() {
			fileName := widgets.QFileDialog_GetOpenFileName(nil, "Open Project", ".", "FireFly project (*.ffp)", "", 0)
			if fileName == "" {
				return
			}
			project, err := storage.LoadFile(fileName)
			if err != nil {
				logrus.Error(err)
				return
			}
			editor.New(project, ApplicationCallbacks)
		},
	}
}
