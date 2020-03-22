package editor

import (
	"image/color"

	"github.com/omniskop/firefly/pkg/project"
	"github.com/omniskop/firefly/pkg/storage"
	"github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

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
	if e.stage.selection == nil {
		return
	}
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
	if e.stage.selection == nil {
		return
	}
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

func (e *Editor) CopyAction(bool) {
	if e.stage.selection == nil {
		return
	}
	e.clipboard = e.stage.selection.element.Copy()
	logrus.Info("copied element")
}

func (e *Editor) PasteAction(bool) {
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

func (e *Editor) deleteSelectedElementAction(bool) {
	if e.stage.selection != nil {
		e.stage.removeElement(e.stage.selection)
	}
}

func (e *Editor) SaveAction(bool) {
	//path := widgets.NewQFileDialog(e.window, core.Qt__Dialog)
	path := widgets.QFileDialog_GetSaveFileName(e.window, "Save the Project", "./project.ffp", "", "", 0)
	err := storage.SaveFile(path, e.project)
	if err != nil {
		logrus.Error(err)
	}
}

func (e *Editor) OpenAction(bool) {
	e.applicationCallbacks["open"]()
}

func (e *Editor) mirrorElementAction(bool) {
	if e.stage.selection == nil {
		return
	}
	copiedElement := e.stage.selection.element.Copy()
	copiedElement.Mirror()
	item := e.stage.addElement(copiedElement)
	item.selectElement()
}
