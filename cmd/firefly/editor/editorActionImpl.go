package editor

import (
	"fmt"
	"image/color"
	"io"
	"math"
	"os"
	"path"
	"path/filepath"
	"reflect"

	"github.com/omniskop/firefly/pkg/project/vectorpath"

	"github.com/omniskop/firefly/cmd/firefly/settings"

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

// PasteAction gets called when the user initiates a paste.
// When there are elements in the clipboard it will put them in the scene according to the current paste mode and
// automatically select them while dropping the previous selection.
func (e *Editor) PasteAction(bool) {
	if len(e.clipboard) == 0 {
		return
	}
	e.stage.selection.clear()

	// To paste all elements at a different location but still relative to each other we first measure the bounding box
	// of all elements.
	bounds := e.clipboard[0].Shape.Bounds() // as a note: can't use infinity rect here or it will cause NaN through calculations
	for _, cbElement := range e.clipboard {
		bounds = bounds.United(cbElement.Shape.Bounds())
	}

	pasteMode := settings.GetString("editor/pasteMode")
	// There are 3 different paste modes:
	// "mouse"
	// 		Pastes the elements as close to the mouse position as possible but without making the elements
	//		leave the stage area. If the copied elements are already partially outside of the stage that will be used
	//		as a limit instead.
	// "needle"
	//      Paste the elements with the earliest at the current needle position and keep the same position axis.
	// "auto"
	// 		Uses "needle" mode while playing and "mouse" when paused.

	// for the mouse mode we will now calculate the starting position for the paste
	var newBoundsLocation vectorpath.Point
	if pasteMode == "mouse" || (pasteMode == "auto" && !e.playing) {
		// Get mouse position in scene coordinates.
		// This will even give correct coordinates when the mouse is outside of the window.
		mousePos := vpPoint(e.stage.MapToScene(e.stage.MapFromGlobal(gui.QCursor_Pos())))

		// Calculate the ideal position for the element's bounds in regards to the cursor
		newBoundsLocation = vectorpath.Point{P: mousePos.P - bounds.Dimensions.P*0.5, T: mousePos.T - bounds.Dimensions.T*0.5}
		// Make sure that this is not outside of the scene or at least not further outside than the copied elements.
		// Left Bounds
		newBoundsLocation.P = math.Max(newBoundsLocation.P, math.Min(bounds.Location.P, 0))
		// Right Bounds
		newBoundsLocation.P = math.Min(newBoundsLocation.P+bounds.Dimensions.P, math.Max(bounds.End().P, 1)) - bounds.Dimensions.P
	}

	for _, cbElement := range e.clipboard {
		element := cbElement.Copy()
		origin := element.Shape.Origin()

		if pasteMode == "mouse" || (pasteMode == "auto" && !e.playing) {
			// adjust this element's position according to the newBoundsLocation that has been calculated earlier
			origin.T = newBoundsLocation.T + (origin.T - bounds.Location.T)
			origin.P = newBoundsLocation.P + (origin.P - bounds.Location.P)
		} else if pasteMode == "needle" || (pasteMode == "auto" && e.playing) {
			// get calculate the offset of this element to the earliest position and add that to the current time
			origin.T = e.Time() + (origin.T - bounds.Location.T)
		}

		element.Shape.SetOrigin(origin)
		element.ZIndex += zIndexSteps // increase ZIndex by one step
		e.stage.selection.add(e.stage.addElement(element))
	}
	logrus.Info("pasted elements")
}

func (e *Editor) CutAction(bool) {
	if e.stage.selection.isEmpty() {
		return
	}
	e.clipboard = e.stage.selection.copyElements()
	e.stage.removeElements(e.stage.selection.elements)
	logrus.Info("cut elements")
}

func (e *Editor) deleteSelectedElementAction(bool) {
	if !e.stage.selection.isEmpty() {
		e.stage.removeElements(e.stage.selection.elements)
	}
}

func (e *Editor) moveToBottomAction(bool) {
	if e.stage.selection.isEmpty() {
		return
	}
	bounds := e.stage.selection.bounds()
	items := e.stage.getItems(bounds)

	// We will search for the highest ZIndex of any element that is in the selection
	// and the lowest ZIndex of all elements that are not selected.
	// Because all elements that are selected will automatically be in the found items
	// we will figure out both values in one go.
	lowestZIndex := math.Inf(1)   // lowest ZIndex of not selected elements
	highestZIndex := math.Inf(-1) // highest ZIndex of selected elements

	for _, item := range items {
		if e.stage.selection.contains(item) {
			if item.element.ZIndex > highestZIndex {
				highestZIndex = item.element.ZIndex
			}
		} else if item.element.ZIndex < lowestZIndex {
			lowestZIndex = item.element.ZIndex
		}
	}
	if lowestZIndex == math.Inf(1) || highestZIndex == math.Inf(-1) {
		// no elements in the bounds
		return
	}

	logrus.WithFields(logrus.Fields{"background lowest": lowestZIndex, "selection highest": highestZIndex}).Info("move selection to bottom")

	for _, item := range e.stage.selection.elements {
		diff := highestZIndex - item.element.ZIndex
		item.element.ZIndex = (lowestZIndex - diff) - zIndexSteps
		item.SetZValue(item.element.ZIndex)
	}
}

func (e *Editor) moveToTopAction(bool) {
	if e.stage.selection.isEmpty() {
		return
	}
	bounds := e.stage.selection.bounds()
	items := e.stage.getItems(bounds)

	// We will search for the lowest ZIndex of any element that is in the selection
	// and the highest ZIndex of all elements that are not selected.
	// Because all elements that are selected will automatically be in the found items
	// we will figure out both values in one go.
	lowestZIndex := math.Inf(1)   // lowest ZIndex of selected elements
	highestZIndex := math.Inf(-1) // highest ZIndex of not selected elements

	for _, item := range items {
		if e.stage.selection.contains(item) {
			if item.element.ZIndex < lowestZIndex {
				lowestZIndex = item.element.ZIndex
			}
		} else if item.element.ZIndex > highestZIndex {
			highestZIndex = item.element.ZIndex
		}
	}
	if lowestZIndex == math.Inf(1) || highestZIndex == math.Inf(-1) {
		// no elements in the bounds
		return
	}

	logrus.WithFields(logrus.Fields{"selection lowest": lowestZIndex, "background highest": highestZIndex}).Info("move selection to bottom")

	for _, item := range e.stage.selection.elements {
		diff := item.element.ZIndex - lowestZIndex
		item.element.ZIndex = (highestZIndex + diff) + zIndexSteps
		item.SetZValue(item.element.ZIndex)
	}
}

func (e *Editor) SaveAction(bool) {
	if e.options.SaveLocation == "" {
		e.SaveAsAction(false)
		return
	}
	err := storage.SaveFile(e.options.SaveLocation, e.project)
	if err != nil {
		logrus.Error(err)
	}
	if e.options.CopyAudioOnSave {
		e.options.CopyAudioOnSave = false
		audioFileName := path.Join(
			filepath.Dir(e.options.SaveLocation),
			fmt.Sprintf("%s - %s%s", e.project.Audio.Title, e.project.Audio.Author, path.Ext(e.player.mediaPath)),
		)

		err := copyFile(audioFileName, e.player.mediaPath)
		if err != nil {
			logrus.Errorf("copy audio file: %w", err)
		}
	}
}

func (e *Editor) SaveAsAction(bool) {
	//path := widgets.NewQFileDialog(e.window, core.Qt__Dialog)
	savePath := widgets.QFileDialog_GetSaveFileName(e.window, "Save the Project", "./project.ffp", "", "", 0)
	e.options.SaveLocation = savePath
	e.SaveAction(false)
}

func (e *Editor) OpenAction(bool) {
	e.applicationCallbacks["open"]()
}

func (e *Editor) OpenSettingsAction(bool) {
	e.applicationCallbacks["openSettings"]()
}

func (e *Editor) mirrorElementAction(bool) {
	if e.stage.selection.isEmpty() {
		return
	}

	copiedElements := e.stage.selection.copyElements()

	if gui.QGuiApplication_KeyboardModifiers()&core.Qt__AltModifier != 0 {
		e.stage.removeElements(e.stage.selection.elements)
	}

	e.stage.selection.clear()
	for _, element := range copiedElements {
		element.Mirror()
		e.stage.selection.add(e.stage.addElement(element))
	}
}

func copyFile(dst, src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil && !os.IsExist(err) {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
