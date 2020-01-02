package editor

import (
	"github.com/omniskop/firefly/cmd/firefly/audio"
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
	window      *widgets.QMainWindow
	project     *project.Project
	stage       *stage
	player      audio.Player
	playing     bool
	updateTimer *core.QTimer
	userActions editorActions
}

func New(proj *project.Project) *Editor {
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
		window:      window,
		project:     proj,
		stage:       nil,
		player:      player,
		playing:     false,
		updateTimer: timer,
		userActions: newEditorActions(),
	}
	edit.userActions.ConnectToEditor(edit)
	edit.stage = newStage(edit, &proj.Scene, proj.Duration)
	window.SetCentralWidget(edit.stage)
	window.AddToolBar(core.Qt__TopToolBarArea, buildEditorToolbar(edit.userActions))

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

func (e *Editor) ToolbarElementAction(checked bool) {
	if e.userActions.group.CheckedAction().Pointer() == e.userActions.cursor.Pointer() {
		e.stage.SetCursor(gui.NewQCursor2(core.Qt__ArrowCursor))
	} else {
		e.stage.SetCursor(gui.NewQCursor2(core.Qt__CrossCursor))
	}
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
	}
}

func (e *Editor) KeyReleaseEvent(event *gui.QKeyEvent) {

}
