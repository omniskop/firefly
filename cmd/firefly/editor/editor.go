package editor

import (
	"image/color"

	"github.com/omniskop/firefly/pkg/storage"

	"github.com/omniskop/firefly/cmd/firefly/audio"
	"github.com/omniskop/firefly/pkg/project"
	"github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const editorViewWidth = 1
const verticalTimeAxis = true

var noPen = gui.NewQPen2(core.Qt__NoPen)

type Editor struct {
	applicationCallbacks map[string]func()
	window               *widgets.QMainWindow
	project              *project.Project
	stage                *stage
	player               audio.Player
	playing              bool
	updateTimer          *core.QTimer
	userActions          *editorActions

	clipboard *project.Element
}

func New(proj *project.Project, applicationCallbacks map[string]func()) *Editor {
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(300, 200)
	window.SetWindowTitle("FireFly Editor")

	player, err := audio.Open(proj.Audio)
	if err != nil {
		logrus.Error(err)
	}

	// Setup update loop
	timer := core.NewQTimer(window)
	timer.SetInterval(1000 / 60)

	edit := &Editor{
		applicationCallbacks: applicationCallbacks,
		window:               window,
		project:              proj,
		stage:                nil,
		player:               player,
		playing:              false,
		updateTimer:          timer,
		userActions:          newEditorActions(),
	}
	edit.userActions.connectToEditor(edit)
	edit.stage = newStage(edit, &proj.Scene, proj.Duration)
	window.SetCentralWidget(edit.stage)
	window.AddToolBar(core.Qt__TopToolBarArea, edit.userActions.buildToolbar())
	window.SetMenuBar(edit.userActions.buildMenuBar())

	window.ConnectKeyPressEvent(edit.KeyPressEvent)
	window.ConnectKeyReleaseEvent(edit.KeyReleaseEvent)

	window.ConnectWheelEvent(func(event *gui.QWheelEvent) {
		event.Ignore()
		window.WheelEventDefault(event)
	})

	window.Show()
	edit.updateTimer.ConnectTimeout(edit.UpdateTick)
	edit.updateTimer.Start2()

	return edit
}

func (e *Editor) UpdateTick() {
	if e.playing {
		audioTime := e.player.Time()
		// fmt.Println("audio time: ", audioTime)
		e.stage.setTime(audioTime)
	}
}

func (e *Editor) UpdateScrollPosition(float64) {

}

func (e *Editor) Time() float64 {
	// TODO: implement AudioOffset
	return e.player.Time()
}

func (e *Editor) SetTime(t float64) {
	// TODO: implement AudioOffset
	e.player.SetTime(t)
	//e.stage.setTime(t)
}

// elementSelected will be called by the stage to notify the editor about an element getting selected
func (e *Editor) elementSelected(item *elementGraphicsItem) {
	if item == nil {
		return
	}
	// update pattern toolbar
	switch e.stage.selection.element.Pattern.(type) {
	case *project.SolidColor:
		e.userActions.solidColor.SetChecked(true)
		e.userActions.colorB.SetDisabled(true)
	case *project.LinearGradient:
		e.userActions.linearGradient.SetChecked(true)
		e.userActions.colorB.SetDisabled(false)
	}
}

func (e *Editor) ToolbarElementAction(checked bool) {
	if e.userActions.toolGroup.CheckedAction().Pointer() == e.userActions.cursor.Pointer() {
		e.stage.SetCursor(gui.NewQCursor2(core.Qt__ArrowCursor))
	} else {
		e.stage.SetCursor(gui.NewQCursor2(core.Qt__CrossCursor))
	}
}

func (e *Editor) ToolbarPatternAction(action *widgets.QAction) {
	if e.stage.selection == nil {
		return
	}
	var col color.Color
	var elementIsSolidColor bool
	switch p := e.stage.selection.element.Pattern.(type) {
	case *project.SolidColor:
		col = p.Color
		elementIsSolidColor = true
	case *project.LinearGradient:
		col = p.Start.Color
	}

	// we only update the pattern when it has changed to prevent endless loops of element updates

	if action.Pointer() == e.userActions.solidColor.Pointer() {
		if elementIsSolidColor {
			return
		}
		e.stage.selection.element.Pattern = project.NewSolidColor(col)
	} else {
		if !elementIsSolidColor {
			return
		}
		e.stage.selection.element.Pattern = project.NewLinearGradient(col, col)
	}
	// TODO: rewrite
	selection := e.stage.selection
	selection.deselectElement()
	selection.selectElement()
	e.stage.selection.updatePattern()
}

func (e *Editor) ToolbarColorAAction(bool) {
	// TODO: rewrite
	var col color.Color
	switch p := e.stage.selection.element.Pattern.(type) {
	case *project.SolidColor:
		col = p.Color
	case *project.LinearGradient:
		col = p.Start.Color
	}

	qcolor := widgets.QColorDialog_GetColor(NewQColorFromColor(col), e.window, "Choose Color", 0)
	if !qcolor.IsValid() { // user canceled dialog
		return
	}
	col = NewColorFromQColor(qcolor)

	switch p := e.stage.selection.element.Pattern.(type) {
	case *project.SolidColor:
		p.Color = col
	case *project.LinearGradient:
		p.Start.Color = col
	}
	e.stage.selection.updatePattern()
}

func (e *Editor) ToolbarColorBAction(bool) {
	// TODO: rewrite
	var col color.Color
	switch p := e.stage.selection.element.Pattern.(type) {
	case *project.SolidColor:
		col = p.Color
	case *project.LinearGradient:
		col = p.Stop.Color
	}

	qcolor := widgets.QColorDialog_GetColor(NewQColorFromColor(col), e.window, "Choose Color", 0)
	if !qcolor.IsValid() { // user canceled dialog
		return
	}
	col = NewColorFromQColor(qcolor)

	switch p := e.stage.selection.element.Pattern.(type) {
	case *project.SolidColor:
		p.Color = col
	case *project.LinearGradient:
		p.Stop.Color = col
	}
	e.stage.selection.updatePattern()
}

func (e *Editor) Copy(bool) {
	if e.stage.selection == nil {
		return
	}
	e.clipboard = e.stage.selection.element.Copy()
	logrus.Info("copied element")
}

func (e *Editor) Paste(bool) {
	if e.clipboard == nil {
		return
	}
	element := e.clipboard.Copy()
	origin := element.Shape.Origin()
	origin.T = e.Time()
	element.Shape.SetOrigin(origin)
	item := e.stage.addElement(element)
	item.selectElement()
	logrus.Info("pasted element")
}

func (e *Editor) Save(bool) {
	//path := widgets.NewQFileDialog(e.window, core.Qt__Dialog)
	path := widgets.QFileDialog_GetSaveFileName(e.window, "Save the Project", "./project.ffp", "", "", 0)
	err := storage.SaveFile(path, e.project)
	if err != nil {
		logrus.Error(err)
	}
}

func (e *Editor) Open(bool) {
	e.applicationCallbacks["open"]()
}

func (e *Editor) KeyPressEvent(event *gui.QKeyEvent) {
	switch core.Qt__Key(event.Key()) {
	case core.Qt__Key_Space:
		logrus.Info("Play/Pause")
		if e.playing {
			e.playing = false
			e.player.Pause()
		} else {
			e.playing = true
			e.player.Play()
		}
	case core.Qt__Key_Backspace:
		if e.stage.selection != nil {
			e.stage.removeElement(e.stage.selection)
		}
	/*case core.Qt__Key_S:
	err := storage.SaveFile("../project_save.json", e.project)
	if err != nil {
		logrus.Errorf("unable to save: %v\n", err)
	} else {
		logrus.Info("file saved!")
	}
	*/
	case core.Qt__Key_9:
		t := e.stage.time()
		logrus.Debug("time is ", t, " ", e.player.Time())
		e.stage.setTime(t)
	case core.Qt__Key_8:
		player := e.player.(*audio.FilePlayer)
		if player.PlaybackRate() == 0.5 {
			player.SetPlaybackRate(1)
		} else {
			player.SetPlaybackRate(0.5)
		}
	}
}

func (e *Editor) KeyReleaseEvent(event *gui.QKeyEvent) {

}
