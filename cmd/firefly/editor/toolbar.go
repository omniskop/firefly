package editor

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type editorActions struct {
	cursor    *widgets.QAction
	newRect   *widgets.QAction
	newTrapez *widgets.QAction
	group     *widgets.QActionGroup
}

func newEditorActions() editorActions {
	actions := editorActions{}

	actions.cursor = newQActionWithIcon("Move", "assets/images/toolbar cursor.imageset/toolbar cursor.png")

	actions.newRect = widgets.NewQAction2("Create Rectangle", nil)
	actions.newRect.SetCheckable(true)
	actions.newRect.SetIcon(gui.NewQIcon5("assets/images/toolbar new rect.imageset/toolbar new rect.png"))

	actions.newTrapez = widgets.NewQAction2("Create Trapezoid", nil)
	actions.newTrapez.SetCheckable(true)
	actions.newTrapez.SetIcon(gui.NewQIcon5("assets/images/toolbar new trapez.imageset/toolbar new trapez.png"))

	actions.group = widgets.NewQActionGroup(nil)
	//TODO: when qt is updated to >= 5.14 set the ExclusionPolicy of the group to QActionGroup::ExclusiveOptional
	actions.group.AddAction(actions.cursor)
	actions.group.AddAction(actions.newRect)
	actions.group.AddAction(actions.newTrapez)

	return actions
}

func (actions editorActions) ConnectToEditor(e *Editor) {
	e.userActions.cursor.ConnectTriggered(e.ToolbarElementAction)
	e.userActions.newRect.ConnectTriggered(e.ToolbarElementAction)
	e.userActions.newTrapez.ConnectTriggered(e.ToolbarElementAction)
}

func newQActionWithIcon(name string, iconPath string) *widgets.QAction {
	action := widgets.NewQAction2(name, nil)
	action.SetCheckable(true)
	action.SetIcon(gui.NewQIcon5(iconPath))
	return action
}

func buildEditorToolbar(actions editorActions) *widgets.QToolBar {
	bar := widgets.NewQToolBar2(nil)
	bar.SetMovable(false)
	bar.SetIconSize(core.NewQSize2(25, 25))

	bar.AddActions([]*widgets.QAction{
		actions.cursor,
		actions.newRect,
		actions.newTrapez,
	})

	return bar
}
