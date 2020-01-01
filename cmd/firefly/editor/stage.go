package editor

import (
	"fmt"
	"image/color"
	"runtime"

	"github.com/omniskop/firefly/pkg/project"
	"github.com/omniskop/firefly/pkg/scanner"
	"github.com/omniskop/firefly/pkg/streamer"

	"github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const needlePosition = 100

type stage struct {
	*widgets.QGraphicsView
	scene        *widgets.QGraphicsScene
	projectScene *project.Scene
	editor       *Editor
	duration     float64
	selection    *elementGraphicsItem

	scanner     scanner.Scanner
	needleFrame scanner.Frame
	streamer    streamer.Streamer
}

func newStage(editor *Editor, projectScene *project.Scene, duration float64) *stage {
	scene := widgets.NewQGraphicsScene(nil)
	scene.SetSceneRect2(0, 0, editorViewWidth, duration)

	udpWriter, err := streamer.NewUDPWriter("192.168.178.35:20202")
	if err != nil {
		fmt.Println(err)
	}

	s := stage{
		QGraphicsView: widgets.NewQGraphicsView(nil),
		scene:         scene,
		projectScene:  projectScene,
		editor:        editor,
		duration:      duration,
		scanner:       scanner.New(projectScene, 60),
		streamer:      streamer.New(udpWriter),
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
	s.FitInView(core.NewQRectF4(0, 0, editorViewWidth, 50), core.Qt__IgnoreAspectRatio)
	s.updateScale()
	s.SetResizeAnchor(widgets.QGraphicsView__AnchorUnderMouse)

	// on macOS the gestures don't need to be grabbed because the operating system will detect them for us
	// if we still grab them qt would do the detection instead of macOS
	if runtime.GOOS != "darwin" {
		s.GrabGesture(core.Qt__PinchGesture, 0) // tell qt we want to know about pinch gestures
	}

	s.ConnectWheelEvent(s.wheelEvent)
	s.ConnectResizeEvent(s.resizeEvent)
	s.ConnectEvent(s.event)
	s.ConnectEventFilter(s.eventFilter)
	s.ConnectDrawForeground(s.drawForeground)
	scene.ConnectChanged(s.sceneChanged)

	s.ConnectScrollContentsBy(func(dx int, dy int) {
		s.ScrollContentsByDefault(dx, dy)
		s.updateNeedleFrame()
	})

	scene.ConnectMousePressEvent(s.sceneMousePressEvent)

	return &s
}

func (s *stage) createElements() {
	s.scene.Clear()

	s.scene.AddRect2(0, 0, editorViewWidth, s.duration, gui.NewQPen2(core.Qt__NoPen), gui.NewQBrush3(gui.NewQColor3(0, 0, 255, 20), core.Qt__SolidPattern))

	for i := range s.projectScene.Elements {
		s.scene.AddItem(newElementGraphicsItem(s, &s.projectScene.Elements[i]))
	}
}

func (s *stage) scaleScene(factor float64) {
	if verticalTimeAxis {
		s.Scale(1, factor)
	} else {
		s.Scale(factor, 1)
	}
	s.updateScale()

	s.setTime(s.editor.Time())
}

func (s *stage) updateScale() {
	physicalZero := s.MapFromScene(core.NewQPointF())
	if verticalTimeAxis {
		point := s.MapToScene(core.NewQPoint2(0, physicalZero.Y()-needlePosition))
		s.scene.SetSceneRect2(0, point.Y(), editorViewWidth, s.duration)
	} else {
		point := s.MapToScene(core.NewQPoint2(physicalZero.X()-needlePosition, 0))
		s.scene.SetSceneRect2(point.Y(), 0, s.duration, editorViewWidth)
	}
}

func (s *stage) time() float64 {
	needleInScene := s.MapToScene(core.NewQPoint2(needlePosition, needlePosition))
	if verticalTimeAxis {
		return needleInScene.Y()
	} else {
		return needleInScene.X()
	}
}

func (s *stage) setTime(t float64) {
	s.scrollSceneToLogical(core.NewQPointF3(t, t), core.NewQPoint2(needlePosition, needlePosition))
}

func (s *stage) updateNeedleFrame() {
	s.needleFrame = s.scanner.Scan(s.time())
	s.streamer.Stream(s.needleFrame)
}

func (s *stage) scrollSceneToLogical(scenePoint *core.QPointF, viewportPoint *core.QPoint) {
	viewPoint := s.Matrix().Map4(scenePoint)

	if verticalTimeAxis {
		s.VerticalScrollBar().SetValue(int(viewPoint.Y()) - viewportPoint.Y())
	} else {
		s.HorizontalScrollBar().SetValue(int(viewPoint.X()) - viewportPoint.X())
	}
}

func (s *stage) sceneChanged([]*core.QRectF) {
	s.updateNeedleFrame()
}

func (s *stage) elementSelected(item *elementGraphicsItem) {
	logrus.Trace("editor element selected")
	if s.selection != item {
		if s.selection != nil {
			logrus.Trace("editor called deselectElement")
			s.selection.deselectElement()
		}
		s.selection = item
		logrus.WithField("item", item).Trace("editor selection changed")
	}
}

func (s *stage) sceneMousePressEvent(event *widgets.QGraphicsSceneMouseEvent) {
	// The builtin selection mechanism has some unwanted side effects that resulted in the need to implement my own.
	// The items themselves will know when they get clicked but I don't know when the used click on the background.
	// The best solution to this problem would probably be to create an item that fills the whole scene and that would
	// received mouse press events when no other item got hit.
	// I tried to fully reimplement the mouse press event without calling the default implementation so that it would
	// only be required to find the clicked item once but that was more complicated than I thought because
	// it doesn't seems possible to use the grabMouse mechanism of qt and I would also need to reimplement that.
	if s.selection != nil {
		hitItem := s.scene.ItemAt(event.ScenePos(), s.ViewportTransform())
		if hitItem == nil {
			s.selection.deselectElement()
			s.selection = nil
		}
	}
	event.Ignore()
	s.scene.MousePressEventDefault(event)
}

func (s *stage) wheelEvent(event *gui.QWheelEvent) {
	// if this is a wheel event and alt is held we want to scale the viewport
	if (event.Modifiers()&core.Qt__AltModifier != 0) && event.Type() == core.QEvent__Wheel {
		// TODO: Use AngleDelta on platforms that don't support PixelDelta (See: https://doc.qt.io/qt-5/qwheelevent.html#pixelDelta)
		deltaY := float64(event.PixelDelta().Y())
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
	s.editor.SetTime(s.time())
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

	// calculate the scaling og the window and then apply the same to the view
	s.Scale(
		float64(event.Size().Width())/float64(event.OldSize().Width()),
		float64(event.Size().Height())/float64(event.OldSize().Height()),
	)
	s.updateScale()
}

func (s *stage) drawForeground(painter *gui.QPainter, rect *core.QRectF) {
	var needleStart *core.QPointF
	var needleStop *core.QPointF
	var gradientStart *core.QPointF
	if verticalTimeAxis {
		needleStart = s.MapToScene(core.NewQPoint2(0, needlePosition))
		needleStop = core.NewQPointF3(editorViewWidth, needleStart.Y())
		gradientStart = s.MapToScene(core.NewQPoint2(0, needlePosition-20))
	} else {
		needleStart = s.MapToScene(core.NewQPoint2(needlePosition, 0))
		needleStop = core.NewQPointF3(needleStart.X(), editorViewWidth)
		gradientStart = s.MapToScene(core.NewQPoint2(needlePosition-20, 0))
	}

	// === draw a line at the position of the needle
	pen := gui.NewQPen3(gui.NewQColor3(255, 255, 255, 255))
	pen.SetWidth(2)
	pen.SetCosmetic(true)

	painter.SetPen(pen)
	painter.DrawLine(core.NewQLineF2(needleStart, needleStop))

	// === draw preview of pixels
	painter.SetPen(noPen)
	//fill := gui.NewQBrush3(gui.NewQColor3(0, 0, 0, 255), core.Qt__SolidPattern)
	pixelWidth := editorViewWidth / float64(len(s.needleFrame.Pixel))
	gradient := gui.NewQLinearGradient2(gradientStart, needleStart)
	for i, pixel := range s.needleFrame.Pixel {
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
}
