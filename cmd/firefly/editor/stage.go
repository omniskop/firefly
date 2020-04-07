package editor

import (
	"fmt"
	"image/color"
	"runtime"
	"unsafe"

	"github.com/omniskop/firefly/pkg/project"
	"github.com/omniskop/firefly/pkg/project/vectorpath"
	"github.com/omniskop/firefly/pkg/scanner"
	"github.com/omniskop/firefly/pkg/streamer"

	"github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const zIndexSteps = 0.1
const infoSideStripe = 150
const audioSideStripe = 150

type stage struct {
	*widgets.QGraphicsView
	scene        *widgets.QGraphicsScene
	projectScene *project.Scene
	editor       *Editor
	duration     float64
	selection    elementList
	items        map[unsafe.Pointer]*elementGraphicsItem

	creationElement *elementGraphicsItem
	creationStart   vectorpath.Point
	newZIndex       float64

	needlePosition int
	needlePipeline *streamer.Pipeline

	nextNonUserScrollEvents uint

	hideElements    bool
	debugShowBounds bool
	debugShowZIndex bool
}

func newStage(editor *Editor, projectScene *project.Scene, duration float64) *stage {
	scene := widgets.NewQGraphicsScene(nil)
	scene.SetSceneRect2(0, 0, editorViewWidth, duration)
	scene.SetBackgroundBrush(gui.NewQBrush3(gui.NewQColor3(14, 15, 16, 255), core.Qt__SolidPattern))

	/*udpWriter, err := streamer.NewUDPWriter("192.168.178.35:20202")
	if err != nil {
		logrus.Error(err)
	}*/

	s := stage{
		QGraphicsView: widgets.NewQGraphicsView(nil),
		scene:         scene,
		projectScene:  projectScene,
		editor:        editor,
		duration:      duration,
		needlePipeline: streamer.NewPipeline(scanner.New(projectScene, 60), streamer.New(nil)),
		selection:      elementList{onChange: editor.selectionChanged},
		items:          make(map[unsafe.Pointer]*elementGraphicsItem),
	}

	s.SetObjectName("mainEditorView")
	s.SetScene(scene)
	s.createElements()
	s.SetRenderHints(gui.QPainter__Antialiasing | gui.QPainter__HighQualityAntialiasing | gui.QPainter__SmoothPixmapTransform)
	// the default viewport-update-mode caused graphical glitches with the macOS scroll bar
	// TODO: check if this is still the case
	s.SetViewportUpdateMode(widgets.QGraphicsView__FullViewportUpdate)
	s.SetVerticalScrollBarPolicy(core.Qt__ScrollBarAlwaysOn)
	s.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	s.SetDragMode(widgets.QGraphicsView__RubberBandDrag)
	s.SetRubberBandSelectionMode(core.Qt__IntersectsItemShape)
	s.updateSize()
	s.updateScale()
	s.showTimeSpan(10)

	// on macOS the gestures don't need to be grabbed because the operating system will detect them for us
	// if we still grab them qt would do the detection instead of macOS
	if runtime.GOOS != "darwin" {
		s.GrabGesture(core.Qt__PinchGesture, 0) // tell qt we want to know about pinch gestures
	}

	s.ConnectWheelEvent(s.wheelEvent)
	s.ConnectResizeEvent(s.resizeEvent)
	s.ConnectMouseReleaseEvent(s.viewMouseReleaseEvent)
	s.ConnectKeyPressEvent(s.keyPressEvent)
	s.ConnectEvent(s.event)
	s.ConnectEventFilter(s.eventFilter)
	s.ConnectDrawBackground(s.drawBackground)
	s.ConnectDrawForeground(s.drawForeground)
	s.ConnectScrollContentsBy(s.scrollContentsByEvent)
	s.scene.ConnectChanged(s.sceneChanged)

	scene.ConnectMousePressEvent(s.sceneMousePressEvent)
	scene.ConnectMouseReleaseEvent(s.sceneMouseReleaseEvent)
	scene.ConnectMouseMoveEvent(s.sceneMouseMoveEvent)

	return &s
}

func (s *stage) createElements() {
	s.items = make(map[unsafe.Pointer]*elementGraphicsItem)
	s.scene.Clear()

	rect := widgets.NewQGraphicsRectItem3(
		0,
		0,
		editorViewWidth,
		s.duration,
		nil,
	)
	rect.SetPen(gui.NewQPen2(core.Qt__NoPen))
	rect.SetBrush(gui.NewQBrush2(core.Qt__NoBrush))
	rect.ConnectMousePressEvent(func(event *widgets.QGraphicsSceneMouseEvent) {
		if s.creationElement == nil && event.Modifiers()&core.Qt__ShiftModifier == 0 {
			s.selection.clear()
		}
		event.Ignore() // ignore this event so that the qt selection can start on this element
	})
	s.scene.AddItem(rect)

	titleContainer := s.scene.AddRect2(0, 0, 100, -30, noPen, gui.NewQBrush3(gui.NewQColor(), core.Qt__NoBrush))
	titleContainer.SetFlags(widgets.QGraphicsItem__ItemIgnoresTransformations)
	songTitle := widgets.NewQGraphicsSimpleTextItem2(fmt.Sprintf("%s - %s", s.editor.project.Audio.Title, s.editor.project.Audio.Author), titleContainer)
	font := gui.NewQFont2("fat_sans_serif", 25, int(gui.QFont__Bold), false)
	font.InsertSubstitutions("fat_sans_serif", []string{"Montserrat", "Futura", "Arial"})
	songTitle.SetFont(font)
	songTitle.SetFlags(widgets.QGraphicsItem__ItemIgnoresTransformations)
	songTitle.SetBrush(gui.NewQBrush3(gui.NewQColor3(201, 201, 201, 255), core.Qt__SolidPattern))
	songTitle.SetPos2(5, -35)

	for i := range s.projectScene.Elements {
		item := newElementGraphicsItem(s, s.projectScene.Elements[i])
		s.items[item.Pointer()] = item
		s.scene.AddItem(item)
	}
}

func (s *stage) addElement(element *project.Element) *elementGraphicsItem {
	s.projectScene.Elements = append(s.projectScene.Elements, element)
	item := newElementGraphicsItem(s, s.projectScene.Elements[len(s.projectScene.Elements)-1])
	s.items[item.Pointer()] = item
	s.scene.AddItem(item)
	return item
}

func (s *stage) addElements(elements []*project.Element) []*elementGraphicsItem {
	var out = make([]*elementGraphicsItem, len(elements))
	for i, element := range elements {
		out[i] = s.addElement(element)
	}
	return out
	}

func (s *stage) removeElement(item *elementGraphicsItem) {
	s.selection.removeIfFound(item)

	delete(s.items, item.Pointer())
	s.scene.RemoveItem(item)
	for i := range s.projectScene.Elements {
		if s.projectScene.Elements[i] == item.element {
			// this copies the last element of the slice over the removed one and shrinks the slice
			lastIndex := len(s.projectScene.Elements) - 1
			s.projectScene.Elements[i] = s.projectScene.Elements[lastIndex]
			s.projectScene.Elements[lastIndex] = nil
			s.projectScene.Elements = s.projectScene.Elements[:lastIndex]
			return
		}
	}
	logrus.Error("an element that should have been deleted could not be found in the scene")
}

func (s *stage) removeElements(items []*elementGraphicsItem) {
	for _, item := range items {
		s.removeElement(item)
	}
}

func (s *stage) updateNeedlePosition() {
	time := s.time()
	//logrus.Debugf("device pixel ratio is: %f", s.editor.window.WindowHandle().Screen().DevicePixelRatio())
	var ratio = s.editor.window.WindowHandle().Screen().DevicePixelRatio()
	if runtime.GOOS == "darwin" {
		ratio = 1 // on macOS the dpi scaling works out of the box
	}
	s.needlePosition = int(100 * ratio)
	s.setTime(time)
	s.updateScale()
	s.InvalidateScene(core.NewQRectF(), widgets.QGraphicsScene__ForegroundLayer)
}

func (s *stage) scaleScene(factor float64) {
	if factor < 0 {
		return // this could sometimes happen and would result in the scene becoming flipped
	}
	if verticalTimeAxis {
		if s.sceneViewport().Height()/factor > s.duration {
			// If the factor would zoom out so far that the scene should show more than the duration of the audio
			// we will limit the factor to the maximum value.
			// This restriction is arbitrary and could potentially be changed to restrict zoom out even further.
			factor = s.sceneViewport().Height() / s.duration
		}
		s.Scale(1, factor)
	} else {
		if s.sceneViewport().Width()/factor > s.duration {
			// If the factor would zoom out so far that the scene should show more than the duration of the audio
			// we will limit the factor to the maximum value.
			// This restriction is arbitrary and could potentially be changed to restrict zoom out even further.
			factor = s.sceneViewport().Width() / s.duration
		}
		s.Scale(factor, 1)
	}
	s.updateScale()

	s.setTime(s.editor.Time())
}

// updateSize will update the stage size and should be called after the windows has been resized
func (s *stage) updateSize() {
	// The following code is inspired by the QGraphicsView::fitInView method
	// It is slightly changed to take care of the different requirements

	// Reset the view scale to 1:1.
	unity := s.Matrix().MapRect(core.NewQRectF4(0, 0, 1, 1))
	if verticalTimeAxis {
		s.Scale(1/unity.Width(), 1)
	} else {
		s.Scale(1, 1/unity.Width())
	}

	// Calculate the x scaling to correctly fit the scene
	if verticalTimeAxis {
		viewRect := s.Viewport().Rect().Adjusted(infoSideStripe, 0, -audioSideStripe, 0)
		ratio := float64(viewRect.Width()) / s.mapToPosition(editorViewWidth)

		s.Scale(ratio, 1)
	} else {
		viewRect := s.Viewport().Rect().Adjusted(0, infoSideStripe, 0, -audioSideStripe)
		ratio := float64(viewRect.Height()) / s.mapToPosition(editorViewWidth)

		s.Scale(1, ratio)
	}

	// The fitInView method would use the centerOn method here

	// NOTE: The QGraphicsView.centerOn method does not work correctly on macOS
	// when the whole scene is visible and thus no scroll bars are required.
	// When the scroll bars in macOS are set to only appear when scrolling they will not take up any space in the widget
	// and instead be overlayed on top of the scene. The centerOn method will still center the scene as if there are
	// scroll bars on the sides and thus it will not be exactly in the center.
	// See https://github.com/qt/qtbase/blob/5.13.2/src/widgets/graphicsview/qgraphicsview.cpp#L398
	// and https://github.com/qt/qtbase/blob/5.13.2/src/plugins/styles/mac/qmacstyle_mac.mm#L2419
}

// updateScale will set the scene rect and should be called after the time scale has been changed
func (s *stage) updateScale() {
	// calculate the scene sizes of the side stripes
	infoSideStripeScene := s.mapToPosition(infoSideStripe)
	audioSideStripeScene := s.mapToPosition(audioSideStripe)

	// adjust the scene rect
	if verticalTimeAxis {
		s.scene.SetSceneRect2(
			-infoSideStripeScene,
			-s.mapToTime(s.needlePosition),
			infoSideStripeScene+editorViewWidth+audioSideStripeScene,
			s.duration+s.mapToTime(s.Viewport().Rect().Height()),
		)
	} else {
		s.scene.SetSceneRect2(
			-s.mapToTime(s.needlePosition),
			-infoSideStripeScene,
			s.duration+s.mapToTime(s.Viewport().Rect().Width()),
			infoSideStripeScene+editorViewWidth+audioSideStripeScene,
		)
	}
	s.nextNonUserScrollEvents++
}

func (s *stage) time() float64 {
	needleInScene := s.MapToScene(core.NewQPoint2(s.needlePosition, s.needlePosition))
	if verticalTimeAxis {
		return needleInScene.Y()
	} else {
		return needleInScene.X()
	}
}

func (s *stage) setTime(t float64) {
	s.scrollSceneToLogical(core.NewQPointF3(t, t), core.NewQPoint2(s.needlePosition, s.needlePosition))
}

func (s *stage) updateNeedleFrame() {
	s.needlePipeline.Update <- s.time()
}

func (s *stage) scrollSceneToLogical(scenePoint *core.QPointF, viewportPoint *core.QPoint) {
	// Inspired by the QGraphicsView.centerOn Method: https://github.com/qt/qtbase/tree/35a461d0261af4178e560df3e3c8fd6fd19bdeb5/src/widgets/graphicsview/qgraphicsview.cpp#L1915

	viewPoint := s.Matrix().Map4(scenePoint)
	s.nextNonUserScrollEvents++

	if verticalTimeAxis {
		s.VerticalScrollBar().SetValue(int(viewPoint.Y()) - viewportPoint.Y())
	} else {
		s.HorizontalScrollBar().SetValue(int(viewPoint.X()) - viewportPoint.X())
	}
}

func (s *stage) sceneViewport() *core.QRectF {
	return s.MapToScene2(s.Viewport().Rect()).BoundingRect()
}

func (s *stage) mapToPosition(diff int) float64 {
	x := s.Matrix().Inverted(nil)
	if verticalTimeAxis {
		return x.Map4(core.NewQPointF3(float64(diff), 0)).X() - x.Map4(core.NewQPointF()).X()
	} else {
		return x.Map4(core.NewQPointF3(0, float64(diff))).Y() - x.Map4(core.NewQPointF()).Y()
	}
}

func (s *stage) mapToTime(diff int) float64 {
	x := s.Matrix().Inverted(nil)
	if verticalTimeAxis {
		return x.Map4(core.NewQPointF3(0, float64(diff))).Y() - x.Map4(core.NewQPointF()).Y()
	} else {
		return x.Map4(core.NewQPointF3(float64(diff), 0)).X() - x.Map4(core.NewQPointF()).X()
	}
}

func (s *stage) redraw() {
	s.scene.Update(s.sceneViewport())
}

func (s *stage) sceneChanged([]*core.QRectF) {
	s.updateNeedleFrame()
}

func (s *stage) elementHasBeenClicked(item *elementGraphicsItem, event *widgets.QGraphicsSceneMouseEvent) {
	if item.parent.creationElement != nil {
		// an element is currently being created and this element should ignore the mouse event
		return
	}

	if event.Modifiers()&core.Qt__ShiftModifier != 0 {
		// shift is held ...
		if s.selection.contains(item) {
			// ... and the element is already selected
			// deselect the element
			s.selection.removeIfFound(item)
		} else {
			// ... and the element is not already selected
			// add the element to the selection
			s.selection.add(item)
		}
	} else if !s.selection.contains(item) {
		// shift is not held and not already selected
		// change the selection to only contain this element
		s.selection.clear()
		s.selection.add(item)
	}
}

func (s *stage) getItems(rect *core.QRectF) []*elementGraphicsItem {
	items := s.scene.Items3(rect, core.Qt__IntersectsItemShape, core.Qt__DescendingOrder, s.Transform())
	var out = make([]*elementGraphicsItem, 0, len(items))
	for _, item := range items {
		if elementItem, ok := s.items[item.Pointer()]; ok {
			out = append(out, elementItem)
		}
		// not all graphicsItems are in the map
	}
	return out
}

func (s *stage) showTimeSpan(timeSpan float64) {
	if timeSpan == 0 {
		logrus.Warn("stage.showTimeSpan has been called with an invalid argument of 0")
		return
	}

	if verticalTimeAxis {
		s.scaleScene(s.sceneViewport().Height() / timeSpan)
	} else {
		s.scaleScene(s.sceneViewport().Width() / timeSpan)
	}
}

func (s *stage) sceneMousePressEvent(event *widgets.QGraphicsSceneMouseEvent) {
	// The builtin selection mechanism has some unwanted side effects that resulted in the need to implement my own.
	// The items themselves will know when they get clicked but I don't know when the user clicks on the background.
	// The best solution to this problem would probably be to create an item that fills the whole scene and that would
	// received mouse press events when no other item got hit.
	// I tried to fully reimplement the mouse press event without calling the default implementation so that it would
	// only be required to find the clicked item once but that was more complicated than I thought because
	// it doesn't seems possible to use the grabMouse mechanism of qt and I would also need to reimplement that.

	if s.editor.userActions.toolGroup.CheckedAction().Pointer() != s.editor.userActions.cursor.Pointer() {
		var elementColor project.Pattern = project.NewSolidColorRGBA(255, 255, 255, 255)
		if !s.selection.isEmpty() {
			// we copy the style of the first selected element
			elementColor = s.selection.elements[0].element.Pattern.Copy()
		}
		s.creationElement = newElementGraphicsItem(s, &project.Element{
			ZIndex:  s.newZIndex,
			Shape:   s.editor.userActions.getSelectedShape(),
			Pattern: elementColor,
		})
		s.newZIndex += zIndexSteps
		s.items[s.creationElement.Pointer()] = s.creationElement
		s.scene.AddItem(s.creationElement)
		s.creationStart = vpPoint(event.ScenePos())
		s.selection.set(s.creationElement)
		logrus.WithField("start", s.creationStart).Debug("a new element is being created")
		// we will not call the default event handler to prevent the rubber band selection from starting
		return
	}

	if !s.selection.isEmpty() {
		hitItem := s.scene.ItemAt(event.ScenePos(), s.ViewportTransform())
		if hitItem == nil {
			s.selection.clear()
		}
	}

	event.Ignore() // ignore means that the event will be handled further
	s.scene.MousePressEventDefault(event)
}

func (s *stage) sceneMouseReleaseEvent(event *widgets.QGraphicsSceneMouseEvent) {
	if s.creationElement != nil {
		// an element is currently being created
		s.projectScene.Elements = append(s.projectScene.Elements, s.creationElement.element)
		// we do this to get the new correct reference to the element in the slice because element is copied
		s.creationElement.element = s.projectScene.Elements[len(s.projectScene.Elements)-1]
		s.creationElement = nil
		s.editor.userActions.cursor.Toggle() // switch the tool back to the standard cursor
	}

	event.Ignore()
	s.scene.MouseReleaseEventDefault(event)
}

func (s *stage) sceneMouseMoveEvent(event *widgets.QGraphicsSceneMouseEvent) {
	if s.creationElement != nil {
		mousePosition := vpPoint(event.ScenePos())
		s.creationElement.element.Shape.SetCreationBounds(s.creationStart, mousePosition.Sub(s.creationStart))
		s.creationElement.updatePath()
	}

	event.Ignore()
	s.scene.MouseMoveEventDefault(event)
}

func (s *stage) wheelEvent(event *gui.QWheelEvent) {
	// if this is a wheel event and a modifier is held we want to scale the viewport
	// on windows the alt modifier does not seem to work so ctrl will be used in there
	// TODO: figure out why that is the case
	// we could potentially make this a setting
	modifier := core.Qt__AltModifier
	if runtime.GOOS == "windows" {
		modifier = core.Qt__ControlModifier
	}
	if (event.Modifiers()&modifier != 0) && event.Type() == core.QEvent__Wheel {
		deltaY := float64(event.PixelDelta().Y())
		// Use AngleDelta on platforms that don't support PixelDelta (See: https://doc.qt.io/qt-5/qwheelevent.html#pixelDelta)
		if deltaY == 0 {
			deltaY = float64(event.AngleDelta().Y()) * 5
		}
		deltaY /= 1000
		s.scaleScene(1 + deltaY)
		return
	}

	// this prevents scrolling on the P Axis
	// we create a new QWheelEvent without the P movement and then call the default wheel event handler
	pixel := event.PixelDelta()
	angle := event.AngleDelta()
	if verticalTimeAxis {
		pixel.SetX(0)
		angle.SetX(0)
	} else {
		pixel.SetY(0)
		angle.SetY(0)
	}
	event = gui.NewQWheelEvent7(event.PosF(), event.GlobalPosF(), pixel, angle, event.Buttons(), event.Modifiers(), event.Phase(), event.Inverted(), event.Source())

	event.Ignore() // TODO: is this necessary?
	s.WheelEventDefault(event)
}

func (s *stage) event(event *core.QEvent) bool {
	if event.Type() == core.QEvent__Gesture {
		// TODO: can this be directly converted to a PinchGesture?
		gestureEvent := widgets.NewQGestureEventFromPointer(event.Pointer())
		gesture := gestureEvent.Gesture(core.Qt__PinchGesture)
		if gesture.Pointer() != nil { // not a PinchGesture
			pinchGesture := widgets.NewQPinchGestureFromPointer(gesture.Pointer())
			s.scaleScene(pinchGesture.ScaleFactor())
			return true
		}
	}

	if event.Type() == core.QEvent__NativeGesture {
		gestureEvent := gui.NewQNativeGestureEventFromPointer(event.Pointer())
		if gestureEvent.GestureType() == core.Qt__ZoomNativeGesture {
			s.scaleScene(1 + gestureEvent.Value())
			return true
		}
	}

	event.Ignore()
	return s.EventDefault(event)
}

func (s *stage) eventFilter(target *core.QObject, event *core.QEvent) bool {
	// TODO: check if the editor is currently playing
	//if event.Type() == core.QEvent__Wheel {
	//	logrus.Debug("event filter")
	//	event.Accept()
	//	return true
	//}
	event.Ignore()
	return false
}

func (s *stage) resizeEvent(event *gui.QResizeEvent) {
	if event.OldSize().Width() == -1 {
		// ignore if this is the first event at the start
		return
	}
	s.nextNonUserScrollEvents++

	// calculate the scaling of the window and then apply the same to the view
	s.Scale(
		float64(event.Size().Width())/float64(event.OldSize().Width()),
		float64(event.Size().Height())/float64(event.OldSize().Height()),
	)
	s.setTime(s.editor.Time())
	s.updateSize()
	s.updateScale()
}

func (s *stage) scrollContentsByEvent(dx int, dy int) {
	s.ScrollContentsByDefault(dx, dy)
	if dy == 0 {
		return
	}
	s.updateNeedleFrame()

	if s.nextNonUserScrollEvents == 0 {
		// do not update the time of the editor because it probably was the editor itself who triggered
		// the scroll event
		s.editor.SetTime(s.time())
	} else {
		// TODO: this whole thing is an ugly hack and needs to be properly redone
		if s.nextNonUserScrollEvents > 2 {
			s.nextNonUserScrollEvents = 2
		}
		s.nextNonUserScrollEvents--
	}
}

func (s *stage) viewMouseReleaseEvent(event *gui.QMouseEvent) {
	if !s.RubberBandRect().IsNull() && s.creationElement == nil {
		// get the selection box in scene coordinates
		rect := s.MapToScene2(s.RubberBandRect()).BoundingRect()
		// get the items that intersect with the box
		items := s.getItems(rect)
		if gui.QGuiApplication_KeyboardModifiers()&core.Qt__ShiftModifier == 0 {
			s.selection.clear()
		}
		for _, item := range items {
			s.selection.add(item)
		}
	}

	event.Ignore()
	s.MouseReleaseEventDefault(event)
}

func (s *stage) keyPressEvent(event *gui.QKeyEvent) {
	var pixelChange = 1
	if event.Modifiers()&core.Qt__ShiftModifier != 0 {
		pixelChange = 10
	}
	switch core.Qt__Key(event.Key()) {
	case core.Qt__Key_Up:
		change := s.MapToScene6(0, 0, pixelChange, pixelChange).BoundingRect()
		if !s.selection.isEmpty() {
			// we only move the first element. The others will automatically be updated too
			s.selection.elements[0].MoveBy(0, -change.Height())
		}
	case core.Qt__Key_Right:
		change := s.MapToScene6(0, 0, pixelChange, pixelChange).BoundingRect()
		if !s.selection.isEmpty() {
			// we only move the first element. The others will automatically be updated too
			s.selection.elements[0].MoveBy(change.Width(), 0)
		}
	case core.Qt__Key_Down:
		change := s.MapToScene6(0, 0, pixelChange, pixelChange).BoundingRect()
		if !s.selection.isEmpty() {
			// we only move the first element. The others will automatically be updated too
			s.selection.elements[0].MoveBy(0, change.Height())
		}
	case core.Qt__Key_Left:
		change := s.MapToScene6(0, 0, pixelChange, pixelChange).BoundingRect()
		if !s.selection.isEmpty() {
			// we only move the first element. The others will automatically be updated too
			s.selection.elements[0].MoveBy(-change.Width(), 0)
		}
	default:
		event.Ignore()
		s.KeyPressEventDefault(event)
	}
}

func (s *stage) drawBackground(painter *gui.QPainter, rect *core.QRectF) {
	// draw background color
	painter.SetPen(noPen)
	painter.SetBrush(NewQBrushFromRGBA(32, 34, 37, 255))
	painter.DrawRect(core.NewQRectF4(rect.Left(), rect.Top(), s.mapToPosition(infoSideStripe), rect.Height()))
	k := s.mapToPosition(audioSideStripe)
	painter.DrawRect(core.NewQRectF4(rect.Right()-k, rect.Top(), k, rect.Height()))

	// draw stage background
	painter.SetBrush(NewQBrushFromRGBA(20, 22, 25, 255))
	painter.DrawRect(core.NewQRectF4(0, rect.Top(), editorViewWidth, rect.Height()))

	// draw guidelines
	pen := gui.NewQPen3(gui.NewQColor3(82, 84, 87, 255))
	pen.SetCosmetic(true)
	painter.SetPen(pen)
	painter.DrawLine(core.NewQLineF3(editorViewWidth*0.25, rect.Top(), editorViewWidth*0.25, rect.Bottom()))
	painter.DrawLine(core.NewQLineF3(editorViewWidth*0.75, rect.Top(), editorViewWidth*0.75, rect.Bottom()))
	painter.DrawLine(core.NewQLineF3(editorViewWidth*0.333, rect.Top(), editorViewWidth*0.333, rect.Bottom()))
	painter.DrawLine(core.NewQLineF3(editorViewWidth*0.666, rect.Top(), editorViewWidth*0.666, rect.Bottom()))
	painter.DrawLine(core.NewQLineF3(editorViewWidth*0.5, rect.Top(), editorViewWidth*0.5, rect.Bottom()))

	// draw start and stop lines
	painter.DrawLine3(0, 0, editorViewWidth, 0)
	painter.DrawLine(core.NewQLineF3(0, s.duration, editorViewWidth, s.duration))

	// draw stage shadow
	painter.SetPen(noPen)
	shadowRect := s.MapToScene6(infoSideStripe-30, 0, 30, s.Viewport().Height()).BoundingRect()
	gradient := gui.NewQLinearGradient2(shadowRect.TopLeft(), shadowRect.TopRight())
	gradient.SetColorAt(0, gui.NewQColor3(0, 0, 0, 0))
	gradient.SetColorAt(1, gui.NewQColor3(0, 0, 0, 70))
	painter.SetBrush(gui.NewQBrush10(gradient))
	painter.DrawRect(shadowRect)

	shadowRect = s.MapToScene6(s.Viewport().Width()-audioSideStripe, 0, 30, s.Viewport().Height()).BoundingRect()
	gradient.SetStart(shadowRect.TopRight())
	gradient.SetFinalStop(shadowRect.TopLeft())
	painter.SetBrush(gui.NewQBrush10(gradient))
	painter.DrawRect(shadowRect)
}

func (s *stage) drawForeground(painter *gui.QPainter, rect *core.QRectF) {
	// === draw the bounding boxes of the selected elements
	if s.debugShowBounds && !s.selection.isEmpty() {
		for _, item := range s.selection.elements {
		painter.SetPen(noPen)
		painter.SetBrush(gui.NewQBrush3(gui.NewQColor3(0, 0, 255, 100), core.Qt__SolidPattern))
			bounds := item.BoundingRect()
			bounds.Translate2(item.Pos())
		painter.DrawRect(bounds)

		pen := gui.NewQPen4(
			gui.NewQBrush3(gui.NewQColor3(0, 255, 255, 255), core.Qt__SolidPattern),
			0,
			core.Qt__SolidLine,
			core.Qt__FlatCap,
			core.Qt__BevelJoin,
		)
		pen.SetCosmetic(true)
		painter.SetPen(pen)
		painter.SetBrush(gui.NewQBrush3(gui.NewQColor3(0, 0, 0, 0), core.Qt__SolidPattern))

			myBounds := item.element.Shape.Bounds()
		painter.DrawRect(core.NewQRectF4(myBounds.Location.P, myBounds.Location.T, myBounds.Dimensions.P, myBounds.Dimensions.T))
	}
	}

	// === draw ZIndex indicators of selected elements
	if s.debugShowZIndex && !s.selection.isEmpty() {
		painter.SetPen(noPen)
		painter.SetBrush(gui.NewQBrush3(gui.NewQColor3(0, 0, 0, 255), core.Qt__SolidPattern))
		for _, item := range s.selection.elements {
			painter.DrawRect(core.NewQRectF4(item.X()+0.02*item.element.ZIndex*10, item.Y(), 0.02, 0.5))
		}
	}

	// === draw a giant rectangle to hide all elements
	// TODO: could this be implemented in a better way?
	if s.hideElements {
		painter.SetPen(noPen)
		painter.SetBrush(NewQBrushFromRGBA(20, 22, 25, 255))
		painter.DrawRect(core.NewQRectF4(0, rect.Top(), editorViewWidth, rect.Height()))
	}

	// === draw the needle and scanner preview
	var needleStart *core.QPointF
	var needleStop *core.QPointF
	var gradientStart *core.QPointF
	if verticalTimeAxis {
		needleStart = s.MapToScene(core.NewQPoint2(0, s.needlePosition))
		needleStop = core.NewQPointF3(rect.Right(), needleStart.Y())
		gradientStart = s.MapToScene(core.NewQPoint2(0, s.needlePosition-30))
	} else {
		needleStart = s.MapToScene(core.NewQPoint2(s.needlePosition, 0))
		needleStop = core.NewQPointF3(rect.Bottom(), editorViewWidth)
		gradientStart = s.MapToScene(core.NewQPoint2(s.needlePosition-30, 0))
	}

	// scanner preview
	painter.SetPen(noPen)
	//fill := gui.NewQBrush3(gui.NewQColor3(0, 0, 0, 255), core.Qt__SolidPattern)
	pixelWidth := editorViewWidth / float64(len(s.needlePipeline.LastFrame.Pixel))
	gradient := gui.NewQLinearGradient2(gradientStart, needleStart)
	for i, pixel := range s.needlePipeline.LastFrame.Pixel {
		//fill.SetColor(NewQColorFromColor(pixel))
		//painter.SetBrush(fill)
		gradient.SetColorAt(1, NewQColorFromColor(pixel))
		gradient.SetColorAt(.5, NewQColorFromColor(pixel))
		gradient.SetColorAt(0, NewQColorFromColor(color.Transparent))
		painter.SetBrush(gui.NewQBrush10(gradient))
		if verticalTimeAxis {
			painter.DrawRect(core.NewQRectF4(float64(i)*pixelWidth, gradientStart.Y(), pixelWidth, needleStart.Y()-gradientStart.Y()))
		} else {
			painter.DrawRect(core.NewQRectF4(gradientStart.Y(), float64(i)*pixelWidth, needleStart.Y()-gradientStart.Y(), pixelWidth))
		}
	}

	// needle
	pen := gui.NewQPen3(gui.NewQColor3(200, 200, 200, 255))
	pen.SetWidth(2)
	pen.SetCosmetic(true)
	painter.SetPen(pen)
	painter.DrawLine(core.NewQLineF2(needleStart, needleStop))
}
