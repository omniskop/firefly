package editor

import (
	"github.com/omniskop/firefly/pkg/project"
	"github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const editorViewWidth = 1000
const verticalTimeAxis = true

var noPen = gui.NewQPen2(core.Qt__NoPen)

type Editor struct {
	window    *widgets.QMainWindow
	project   *project.Project
	view      *widgets.QGraphicsView
	scene     *widgets.QGraphicsScene
}

func New(proj *project.Project) *Editor {
	window := createEditorWindow()

	viewObject := window.FindChild("mainEditorView", core.Qt__FindChildrenRecursively)

	if viewObject == nil {
		// this should never happen but if it does it is a reason to panic
		panic("could not find graphics view in editor window")
	}

	// cast the QObject to a QGraphicsView
	view := widgets.NewQGraphicsViewFromPointer(viewObject.Pointer())

	//view.Scene().SetSceneRect2(0, 0, editorViewWidth, proj.Duration)

	edit := &Editor{
		window:  window,
		project: proj,
		view:    view,
		scene:   view.Scene(),
	}

	window.Show()

	return edit
}

func createEditorWindow() *widgets.QMainWindow {
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(300, 200)
	window.SetWindowTitle("FireFly Editor")

	window.SetCentralWidget(buildEditor())

	return window
}

func buildEditor() *widgets.QWidget {
	mainWidget := widgets.NewQWidget(nil, 0)
	mainLayout := widgets.NewQVBoxLayout()
	mainWidget.SetLayout(mainLayout)

	// mainLayout.SetSpacing(0)
	mainLayout.SetContentsMargins(0, 0, 0, 0)

	label := widgets.NewQLabel2("Hallo!", nil, 0)
	mainWidget.Layout().AddWidget(label)

	stageView := buildStage()
	mainWidget.Layout().AddWidget(stageView)

	return mainWidget
}

