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
	selection *elementGraphicsItem
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
	edit.setupScene()
	view.Scene().ConnectMousePressEvent(edit.mousePressEvent)

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

func (e *Editor) setupScene() {
	logrus.WithField("elements", len(e.project.Scene.Elements)).Debug("editor scene setup")
	e.scene.Clear()

	e.view.Scene().AddRect2(0, 0, editorViewWidth, e.project.Duration, gui.NewQPen2(core.Qt__NoPen), gui.NewQBrush3(gui.NewQColor3(0, 0, 255, 20), core.Qt__SolidPattern))

	for i := range e.project.Scene.Elements {
		e.scene.AddItem(newElementGraphicsItem(e, &e.project.Scene.Elements[i]))
	}
}

func (e *Editor) mousePressEvent(event *widgets.QGraphicsSceneMouseEvent) {
	// The builtin selection mechanism has some unwanted side effects that resulted in the need to implement my own.
	// The items themselves will know when they get clicked but I don't know when the used click on the background.
	// The best solution to this problem would probably be to create an item that fills the whole scene and that would
	// received mouse press events when no other item got hit.
	// I tried to fully reimplement the mouse press event without calling the default implementation so that it would
	// only be required to find the clicked item once but that was more complicated than I thought because
	// it doesn't seems possible to use the grabMouse mechanism of qt and I would also need to reimplement that.
	if e.selection != nil {
		hitItem := e.scene.ItemAt(event.ScenePos(), e.view.ViewportTransform())
		if hitItem == nil {
			e.selection.deselect()
			e.selection = nil
		}
	}
	event.Ignore()
	e.scene.MousePressEventDefault(event)
}

func (e *Editor) elementSelected(item *elementGraphicsItem) {
	logrus.Trace("editor element selected")
	if e.selection != item {
		if e.selection != nil {
			logrus.Trace("editor called deselect")
			e.selection.deselect()
		}
		e.selection = item
		logrus.WithField("item", item).Trace("editor selection changed")
	}
}
