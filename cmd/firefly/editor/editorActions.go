package editor

import (
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

	save   *widgets.QAction
	saveAs *widgets.QAction
	open   *widgets.QAction

	copy          *widgets.QAction
	paste         *widgets.QAction
	cut           *widgets.QAction
	mirrorElement *widgets.QAction
	delete        *widgets.QAction

	solidColor     *widgets.QAction
	linearGradient *widgets.QAction
	patternGroup   *widgets.QActionGroup
	colorA         *widgets.QAction
	colorB         *widgets.QAction
}

func newEditorActions() *editorActions {
	var actions = new(editorActions)

	actions.cursor = newCheckableQActionWithIcon("Move", ":assets/images/toolbar cursor.imageset/toolbar cursor.png")
	actions.cursor.SetShortcut(gui.NewQKeySequence3(int(core.Qt__Key_V), 0, 0, 0))
	actions.cursor.SetChecked(true)

	actions.newRect = newCheckableQActionWithIcon("Create Rectangle", ":assets/images/toolbar new rect.imageset/toolbar new rect.png")
	actions.newTrapez = newCheckableQActionWithIcon("Create Trapezoid", ":assets/images/toolbar new trapez.imageset/toolbar new trapez.png")

	actions.toolGroup = widgets.NewQActionGroup(nil)
	//TODO: when qt is updated to >= 5.14 set the ExclusionPolicy of the toolGroup to QActionGroup::ExclusiveOptional
	actions.toolGroup.AddAction(actions.cursor)
	actions.toolGroup.AddAction(actions.newRect)
	actions.toolGroup.AddAction(actions.newTrapez)

	actions.save = newQActionWithIcon("Save", ":assets/images/toolbar save.imageset/toolbar save.png")
	actions.save.SetShortcut(gui.NewQKeySequence5(gui.QKeySequence__Save))
	actions.saveAs = newQActionWithIcon("Save As...", ":assets/images/toolbar save.imageset/toolbar save.png")
	actions.saveAs.SetShortcut(gui.NewQKeySequence5(gui.QKeySequence__SaveAs))
	actions.open = newQActionWithIcon("Open...", ":assets/images/toolbar open.imageset/toolbar open.png")
	actions.open.SetShortcut(gui.NewQKeySequence5(gui.QKeySequence__Open))

	actions.copy = widgets.NewQAction2("Copy", nil)
	actions.copy.SetShortcut(gui.NewQKeySequence5(gui.QKeySequence__Copy))
	actions.paste = widgets.NewQAction2("Paste", nil)
	actions.paste.SetShortcut(gui.NewQKeySequence5(gui.QKeySequence__Paste))
	actions.cut = widgets.NewQAction2("Cut", nil)
	actions.cut.SetShortcut(gui.NewQKeySequence5(gui.QKeySequence__Cut))
	actions.mirrorElement = widgets.NewQAction2("Mirror", nil)
	actions.mirrorElement.SetShortcut(gui.NewQKeySequence2("m", gui.QKeySequence__NativeText))
	actions.delete = widgets.NewQAction2("Delete", nil)
	actions.delete.SetShortcuts([]*gui.QKeySequence{newQKeySequenceFromKeys(core.Qt__Key_Backspace), newQKeySequenceFromKeys(core.Qt__Key_Delete)}) // Qt__KeySequence_Backspace would not work on macOS

	actions.solidColor = newCheckableQActionWithIcon("Solid Color", ":assets/images/toolbar solid color.imageset/toolbar solid color.png")
	actions.solidColor.SetChecked(true)
	actions.linearGradient = newCheckableQActionWithIcon("Linear Gradient", ":assets/images/toolbar linear gradient.imageset/toolbar linear gradient.png")
	actions.patternGroup = widgets.NewQActionGroup(nil)
	actions.patternGroup.AddAction(actions.solidColor)
	actions.patternGroup.AddAction(actions.linearGradient)
	actions.colorA = newQActionWithIcon("Choose Color", ":assets/images/toolbar colorpicker.imageset/toolbar colorpicker.png")
	actions.colorB = newQActionWithIcon("Choose Second Color", ":assets/images/toolbar colorpicker.imageset/toolbar colorpicker.png")
	actions.colorB.SetDisabled(true)

	return actions
}

func (actions *editorActions) connectToEditor(e *Editor) {
	e.userActions.cursor.ConnectTriggered(e.ToolbarElementAction)
	e.userActions.newRect.ConnectTriggered(e.ToolbarElementAction)
	e.userActions.newTrapez.ConnectTriggered(e.ToolbarElementAction)
	e.userActions.save.ConnectTriggered(e.SaveAction)
	e.userActions.saveAs.ConnectTriggered(e.SaveAsAction)
	e.userActions.open.ConnectTriggered(e.OpenAction)
	e.userActions.copy.ConnectTriggered(e.CopyAction)
	e.userActions.paste.ConnectTriggered(e.PasteAction)
	e.userActions.cut.ConnectTriggered(e.CutAction)
	e.userActions.mirrorElement.ConnectTriggered(e.mirrorElementAction)
	e.userActions.delete.ConnectTriggered(e.deleteSelectedElementAction)
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
		actions.saveAs,
	})

	editMenu := menubar.AddMenu2("Edit")
	editMenu.AddActions([]*widgets.QAction{
		actions.delete,
		actions.copy,
		actions.paste,
		actions.cut,
	})
	editMenu.AddSeparator()
	editMenu.AddActions([]*widgets.QAction{
		actions.mirrorElement,
	})

	return menubar
}

func newCheckableQActionWithIcon(name string, iconPath string) *widgets.QAction {
	action := widgets.NewQAction2(name, nil)
	action.SetCheckable(true)
	action.SetIcon(gui.NewQIcon5(iconPath))
	return action
}

func newQActionWithIcon(name string, iconPath string) *widgets.QAction {
	action := widgets.NewQAction2(name, nil)
	action.SetIcon(gui.NewQIcon5(iconPath))
	return action
}

func newQKeySequenceFromKeys(firstKey core.Qt__Key, additional ...core.Qt__Key) *gui.QKeySequence {
	switch len(additional) {
	case 0:
		return gui.NewQKeySequence3(int(firstKey), 0, 0, 0)
	case 1:
		return gui.NewQKeySequence3(int(firstKey), int(additional[0]), 0, 0)
	case 2:
		return gui.NewQKeySequence3(int(firstKey), int(additional[1]), int(additional[2]), 0)
	case 3:
		fallthrough
	default:
		return gui.NewQKeySequence3(int(firstKey), int(additional[1]), int(additional[2]), int(additional[3]))
	}
}
