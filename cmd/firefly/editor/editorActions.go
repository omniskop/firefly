package editor

import (
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/omniskop/firefly/pkg/project/shape"
	"github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type editorActions struct {
	cursor    *widgets.QAction
	newRect   *widgets.QAction
	newTrapez *widgets.QAction
	toolGroup *widgets.QActionGroup

	save  *widgets.QAction
	open  *widgets.QAction
	copy  *widgets.QAction
	paste *widgets.QAction

	solidColor     *widgets.QAction
	linearGradient *widgets.QAction
	patternGroup   *widgets.QActionGroup
	colorA         *widgets.QAction
	colorB         *widgets.QAction
}

func newEditorActions() *editorActions {
	var actions = new(editorActions)

	actions.cursor = newCheckableQActionWithIcon("Move", "assets/images/toolbar cursor.imageset/toolbar cursor.png")
	actions.cursor.SetShortcut(gui.NewQKeySequence3(int(core.Qt__Key_V), 0, 0, 0))
	actions.cursor.SetChecked(true)

	actions.newRect = newCheckableQActionWithIcon("Create Rectangle", "assets/images/toolbar new rect.imageset/toolbar new rect.png")
	actions.newTrapez = newCheckableQActionWithIcon("Create Trapezoid", "assets/images/toolbar new trapez.imageset/toolbar new trapez.png")

	actions.toolGroup = widgets.NewQActionGroup(nil)
	//TODO: when qt is updated to >= 5.14 set the ExclusionPolicy of the toolGroup to QActionGroup::ExclusiveOptional
	actions.toolGroup.AddAction(actions.cursor)
	actions.toolGroup.AddAction(actions.newRect)
	actions.toolGroup.AddAction(actions.newTrapez)

	actions.save = newQActionWithIcon("Save...", "assets/images/toolbar save.imageset/toolbar save.png")
	actions.save.SetShortcut(gui.NewQKeySequence5(gui.QKeySequence__Save))

	actions.open = newQActionWithIcon("Open...", "assets/images/toolbar open.imageset/toolbar open.png")
	actions.open.SetShortcut(gui.NewQKeySequence5(gui.QKeySequence__Open))

	actions.copy = widgets.NewQAction2("Copy", nil)
	actions.copy.SetShortcut(gui.NewQKeySequence5(gui.QKeySequence__Copy))
	actions.paste = widgets.NewQAction2("Paste", nil)
	actions.paste.SetShortcut(gui.NewQKeySequence5(gui.QKeySequence__Paste))

	actions.solidColor = newCheckableQActionWithIcon("Solid Color", "assets/images/toolbar solid color.imageset/toolbar solid color.png")
	actions.solidColor.SetChecked(true)
	actions.linearGradient = newCheckableQActionWithIcon("Linear Gradient", "assets/images/toolbar linear gradient.imageset/toolbar linear gradient.png")
	actions.patternGroup = widgets.NewQActionGroup(nil)
	actions.patternGroup.AddAction(actions.solidColor)
	actions.patternGroup.AddAction(actions.linearGradient)
	actions.colorA = newQActionWithIcon("Choose Color", "assets/images/toolbar colorpicker.imageset/toolbar colorpicker.png")
	actions.colorB = newQActionWithIcon("Choose Second Color", "assets/images/toolbar colorpicker.imageset/toolbar colorpicker.png")
	actions.colorB.SetDisabled(true)

	return actions
}

func (actions *editorActions) connectToEditor(e *Editor) {
	e.userActions.cursor.ConnectTriggered(e.ToolbarElementAction)
	e.userActions.newRect.ConnectTriggered(e.ToolbarElementAction)
	e.userActions.newTrapez.ConnectTriggered(e.ToolbarElementAction)
	e.userActions.save.ConnectTriggered(e.Save)
	e.userActions.open.ConnectTriggered(e.Open)
	e.userActions.copy.ConnectTriggered(e.Copy)
	e.userActions.paste.ConnectTriggered(e.Paste)
	e.userActions.patternGroup.ConnectTriggered(e.ToolbarPatternAction)
	e.userActions.colorA.ConnectTriggered(e.ToolbarColorAAction)
	e.userActions.colorB.ConnectTriggered(e.ToolbarColorBAction)
}

func (actions *editorActions) getSelectedShape() shape.Shape {
	switch actions.toolGroup.CheckedAction().Pointer() {
	case actions.newRect.Pointer():
		return shape.NewEmptyOrthogonalRectangle()
	case actions.newTrapez.Pointer():
		return shape.NewEmptyBentTrapezoid()
	default:
		logrus.Error("a toolbar action is selected that has no known shape that can be created in the stage")
		logrus.Error("the pointer to the action is: ", actions.toolGroup.CheckedAction().Pointer())
		logrus.Errorf("all known action pointers are: \n%v", actions)
		return nil
	}
}

func (actions *editorActions) buildToolbar() *widgets.QToolBar {
	bar := widgets.NewQToolBar2(nil)
	bar.SetMovable(false)
	if runtime.GOOS == "darwin" {
		bar.SetIconSize(core.NewQSize2(25, 25))
	}

	bar.AddActions([]*widgets.QAction{
		actions.cursor,
		actions.newRect,
		actions.newTrapez,
	})
	bar.AddSeparator()
	bar.AddActions([]*widgets.QAction{
		actions.save,
		actions.open,
	})
	bar.AddSeparator()
	bar.AddActions([]*widgets.QAction{
		actions.solidColor,
		actions.linearGradient,
		actions.colorA,
		actions.colorB,
	})

	return bar
}

func (actions *editorActions) buildMenuBar() *widgets.QMenuBar {
	menubar := widgets.NewQMenuBar(nil)

	fileMenu := menubar.AddMenu2("File")
	fileMenu.AddActions([]*widgets.QAction{
		actions.open,
	})
	fileMenu.AddSeparator()
	fileMenu.AddActions([]*widgets.QAction{
		actions.save,
	})

	editMenu := menubar.AddMenu2("Edit")
	editMenu.AddActions([]*widgets.QAction{
		actions.copy,
		actions.paste,
	})

	return menubar
}

func newCheckableQActionWithIcon(name string, iconPath string) *widgets.QAction {
	action := widgets.NewQAction2(name, nil)
	action.SetCheckable(true)
	action.SetIcon(gui.NewQIcon5(path.Join(getPathPrefix(), iconPath)))
	return action
}

func newQActionWithIcon(name string, iconPath string) *widgets.QAction {
	action := widgets.NewQAction2(name, nil)
	action.SetIcon(gui.NewQIcon5(path.Join(getPathPrefix(), iconPath)))
	return action
}

func getPathPrefix() string {
	if runtime.GOOS == "darwin" {
		name, _ := os.Executable()
		return filepath.Dir(name)
		//return filepath.Dir(name) + "/../.."
	}
	return ""
}