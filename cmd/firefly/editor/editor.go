package editor

import (
	"math"
	"reflect"

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
	player               *audioPlayer
	playing              bool
	userActions          *editorActions

	clipboard []*project.Element

	SaveLocation string // the location where the file is saved
}

func New(proj *project.Project, applicationCallbacks map[string]func()) *Editor {
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(300, 200)
	window.SetWindowTitle("Firefly Editor")

	audioPath, err := LocateAudioFile(proj.Audio)
	if err != nil {
		logrus.Error(err)
	}
	player := NewAudioPlayer(audioPath)

	edit := &Editor{
		applicationCallbacks: applicationCallbacks,
		window:               window,
		project:              proj,
		stage:                nil,
		player:               player,
		playing:              false,
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
	gui.NewQWindowFromPointer(window.WindowHandle().Pointer()).ConnectScreenChanged(edit.ScreenChangedEvent)
	edit.stage.updateNeedlePosition() // this needs to be called after the window is shown
	player.onTimeChanged(edit.stage.setTime)

	size := gui.QGuiApplication_PrimaryScreen().AvailableSize()
	size.SetWidth(int(float64(size.Width()) * 0.6))
	size.SetHeight(int(float64(size.Height()) * 0.8))
	window.SetGeometry(widgets.QStyle_AlignedRect(core.Qt__LeftToRight, core.Qt__AlignCenter, size, gui.QGuiApplication_PrimaryScreen().AvailableGeometry()))
	window.Resize(size)

	return edit
}

func (e *Editor) UpdateScrollPosition(float64) {

}

func (e *Editor) Time() float64 {
	// TODO: implement AudioOffset
	return e.player.time()
}

func (e *Editor) SetTime(t float64) {
	// TODO: implement AudioOffset
	e.player.setTime(t)
	//e.stage.setTime(t)
}

// elementSelected will be called by the stage to notify the editor about an element getting selected
func (e *Editor) selectionChanged() {
	e.updateToolbar()
}

// updateToolbar updates the buttons in the toolbar to reflect the current state
func (e *Editor) updateToolbar() {
	if e.stage.selection.isEmpty() {
		return
	}

	// find out if all selected elements have the same pattern type
	var patternType string = reflect.TypeOf(e.stage.selection.elements[0].element.Pattern).String()
	for _, item := range e.stage.selection.elements {
		if reflect.TypeOf(item.element.Pattern).String() != patternType {
			goto patternsOfDifferentType // a different type was found
		}
	}

	// if they are all of the same type, find out which and update the toolbar accordingly
	switch e.stage.selection.elements[0].element.Pattern.(type) {
	case *project.SolidColor:
		e.userActions.solidColor.SetChecked(true)
		e.userActions.colorA.SetDisabled(false)
		e.userActions.colorB.SetDisabled(true)
	case *project.LinearGradient:
		e.userActions.linearGradient.SetChecked(true)
		e.userActions.colorA.SetDisabled(false)
		e.userActions.colorB.SetDisabled(false)
	}

	return
patternsOfDifferentType:
	// the patterns are not the same
	e.userActions.linearGradient.SetChecked(false)
	e.userActions.solidColor.SetChecked(false)
	e.userActions.colorA.SetDisabled(true)
	e.userActions.colorB.SetDisabled(true)
}

func (e *Editor) KeyPressEvent(event *gui.QKeyEvent) {
	switch core.Qt__Key(event.Key()) {
	case core.Qt__Key_Space:
		logrus.Info("Play/Pause")
		if e.playing {
			e.playing = false
			e.player.pause()
		} else {
			e.playing = true
			e.player.play()
		}
	case core.Qt__Key_Minus:
		e.stage.scaleScene(0.9)
	case core.Qt__Key_Plus:
		e.stage.scaleScene(1.1)
	case core.Qt__Key_H:
		e.stage.hideElements = true
		e.stage.redraw()
	/*case core.Qt__Key_S:
	err := storage.SaveFile("../project_save.json", e.project)
	if err != nil {
		logrus.Errorf("unable to save: %v\n", err)
	} else {
		logrus.Info("file saved!")
	}
	*/
	case core.Qt__Key_1:
		t := e.stage.time()
		logrus.Debug("time is ", t, " ", e.player.time())
		e.stage.setTime(t)
	case core.Qt__Key_2:
		e.stage.debugShowBounds = !e.stage.debugShowBounds
		e.stage.redraw()
	case core.Qt__Key_3:
		e.stage.debugShowZIndex = !e.stage.debugShowZIndex
		e.stage.redraw()
	case core.Qt__Key_4:
		pixel := e.stage.needlePipeline.LastFrame.Pixel
		c := pixel[len(pixel)/2]
		r, g, b, a := c.RGBA()
		logrus.Debugf("direct: %d %d %d %d | gamma: %.f %.f %.f %.f",
			r/257, g/257, b/257, a/257,
			math.Pow(float64(r)/0xffff, 2.2)*255,
			math.Pow(float64(g)/0xffff, 2.2)*255,
			math.Pow(float64(b)/0xffff, 2.2)*255,
			math.Pow(float64(a)/0xffff, 2.2)*255)
	case core.Qt__Key_0:
		e.player.SetPlaybackRate(1)
	case core.Qt__Key_9:
		e.player.SetPlaybackRate(0.5)
	case core.Qt__Key_8:
		e.player.SetPlaybackRate(0.25)
	}
}

func (e *Editor) KeyReleaseEvent(event *gui.QKeyEvent) {
	switch core.Qt__Key(event.Key()) {
	case core.Qt__Key_H:
		e.stage.hideElements = false
		e.stage.redraw()
	}
}

func (e *Editor) ScreenChangedEvent(screen *gui.QScreen) {
	e.stage.updateNeedlePosition()
}
