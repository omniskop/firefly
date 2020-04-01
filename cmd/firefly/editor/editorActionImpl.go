package editor

import (
	"image/color"
	"math"
	"reflect"

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
	if e.stage.selection.isEmpty() {
		return
	}

	// we just need a color that can be used and we take one from the first element
	var col color.Color
	switch p := e.stage.selection.elements[0].element.Pattern.(type) {
	case *project.SolidColor:
		col = p.Color
	case *project.LinearGradient:
		col = p.Start.Color
	}

	// now we create the new pattern
	var newPattern project.Pattern
	if e.userActions.solidColor.IsChecked() {
		newPattern = project.NewSolidColor(col)
	} else if e.userActions.linearGradient.IsChecked() {
		newPattern = project.NewLinearGradient(col, col)
	} else {
		return
	}

	newPatternType := reflect.TypeOf(newPattern).String()
	for _, item := range e.stage.selection.elements {
		if reflect.TypeOf(item.element.Pattern).String() != newPatternType {
			item.element.Pattern = newPattern
			item.updatePattern()
		}
	}

	e.updateToolbar()
}

func (e *Editor) ToolbarColorAAction(bool) {
	if e.stage.selection.isEmpty() {
		return
	}

	// get the current color
	// when multiple elements are selected we take the first one
	var col color.Color
	switch p := e.stage.selection.elements[0].element.Pattern.(type) {
	case *project.SolidColor:
		col = p.Color
	case *project.LinearGradient:
		col = p.Start.Color
	}

	qcolor := widgets.QColorDialog_GetColor(NewQColorFromColor(col), e.window, "Choose Color", widgets.QColorDialog__ShowAlphaChannel)
	if !qcolor.IsValid() { // user canceled dialog
		return
	}
	col = NewColorFromQColor(qcolor)

	for _, item := range e.stage.selection.elements {
		switch p := item.element.Pattern.(type) {
		case *project.SolidColor:
			p.Color = col
		case *project.LinearGradient:
			p.Start.Color = col
		}
		item.updatePattern()
	}
}

func (e *Editor) ToolbarColorBAction(bool) {
	if e.stage.selection.isEmpty() {
		return
	}

	// get the current color
	// when multiple elements are selected we take the first one
	var col color.Color
	switch p := e.stage.selection.elements[0].element.Pattern.(type) {
	case *project.SolidColor:
		col = p.Color
	case *project.LinearGradient:
		col = p.Stop.Color
	}

	qcolor := widgets.QColorDialog_GetColor(NewQColorFromColor(col), e.window, "Choose Color", widgets.QColorDialog__ShowAlphaChannel)
	if !qcolor.IsValid() { // user canceled dialog
		return
	}
	col = NewColorFromQColor(qcolor)

	for _, item := range e.stage.selection.elements {
		switch p := item.element.Pattern.(type) {
		case *project.SolidColor:
			p.Color = col
		case *project.LinearGradient:
			p.Stop.Color = col
		}
		item.updatePattern()
	}
}

func (e *Editor) CopyAction(bool) {
	if e.stage.selection.isEmpty() {
		return
	}
	e.clipboard = e.stage.selection.copyElements()
	logrus.Info("copied elements")
}

func (e *Editor) PasteAction(bool) {
	if e.clipboard == nil {
		return
	}
	e.stage.selection.clear()

	// to be able to paste the elements while maintaining relative positioning
	// we need to figure out the earliest origin of the elements in the clipboard
	var earliestPosition = math.Inf(1)
	for _, cbElement := range e.clipboard {
		if cbElement.Shape.Origin().T < earliestPosition {
			earliestPosition = cbElement.Shape.Origin().T
		}
	}

	for _, cbElement := range e.clipboard {
		element := cbElement.Copy()
		origin := element.Shape.Origin()
		// get calculate the offset of this element to the earliest position and add that to the current time
		origin.T = e.Time() + (origin.T - earliestPosition)
		//origin.T = e.Time()
		element.Shape.SetOrigin(origin)
		e.stage.selection.add(e.stage.addElement(element))
	}
	logrus.Info("pasted elements")
}

func (e *Editor) deleteSelectedElementAction(bool) {
	if !e.stage.selection.isEmpty() {
		e.stage.removeElements(e.stage.selection.elements)
	}
}

func (e *Editor) SaveAction(bool) {
	if e.SaveLocation == "" {
		e.SaveAsAction(false)
		return
	}
	err := storage.SaveFile(e.SaveLocation, e.project)
	if err != nil {
		logrus.Error(err)
	}
}

func (e *Editor) SaveAsAction(bool) {
	//path := widgets.NewQFileDialog(e.window, core.Qt__Dialog)
	path := widgets.QFileDialog_GetSaveFileName(e.window, "Save the Project", "./project.ffp", "", "", 0)
	err := storage.SaveFile(path, e.project)
	if err != nil {
		logrus.Error(err)
	}
	e.SaveLocation = path
}

func (e *Editor) OpenAction(bool) {
	e.applicationCallbacks["open"]()
}

func (e *Editor) mirrorElementAction(bool) {
	if e.stage.selection.isEmpty() {
		return
	}

	copiedElements := e.stage.selection.copyElements()
	e.stage.selection.clear()
	for _, element := range copiedElements {
		element.Mirror()
		e.stage.selection.add(e.stage.addElement(element))
	}
}
