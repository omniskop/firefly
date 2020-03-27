package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/uitools"

	"github.com/omniskop/firefly/cmd/firefly/editor"
	"github.com/omniskop/firefly/pkg/storage"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

var errUserAbort = errors.New("user aborted")

func createApplication() *widgets.QApplication {
	core.QCoreApplication_SetAttribute(core.Qt__AA_ShareOpenGLContexts, true)
	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetWindowIcon(gui.NewQIcon5(":assets/images/logo.png"))
	app.SetAttribute(core.Qt__AA_UseHighDpiPixmaps, true)
	return app
}

var ApplicationCallbacks map[string]func()

func init() {
	ApplicationCallbacks = map[string]func(){
		"open": func() {
			_ = openProject()
		},
	}
}

func loadUI(filePath string, parent *widgets.QWidget) (*widgets.QWidget, error) {
	file := core.NewQFile2(filePath)
	ok := file.Open(core.QIODevice__ReadOnly)
	if !ok {
		return nil, fmt.Errorf("can't open file: %s", filePath)
	}
	formWidget := uitools.NewQUiLoader(nil).Load(file, parent)
	file.Close()
	return formWidget, nil
}

func openProject() error {
	fileName := widgets.QFileDialog_GetOpenFileName(nil, "Open Project", ".", "FireFly project (*.ffp)", "", 0)
	if fileName == "" {
		return errUserAbort
	}
	project, err := storage.LoadFile(fileName)
	if err != nil {
		return err
	}
	edit := editor.New(project, ApplicationCallbacks)
	edit.SaveLocation = fileName
	return nil
}
