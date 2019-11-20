package editor

import (
	"runtime"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func buildStage() *widgets.QGraphicsView {
	scene := widgets.NewQGraphicsScene(nil)
	scene.SetSceneRect2(0, 0, editorViewWidth, 300)

	view := widgets.NewQGraphicsView(nil)
	view.SetObjectName("mainEditorView")

	// view.SetViewport(widgets.NewQOpenGLWidget(nil, 0))
	// view.SetViewportUpdateMode(widgets.QGraphicsView__FullViewportUpdate)

	view.SetScene(scene)

	view.SetRenderHint(gui.QPainter__Antialiasing, true)
	view.SetRenderHint(gui.QPainter__HighQualityAntialiasing, true)
	view.SetRenderHint(gui.QPainter__SmoothPixmapTransform, true)

	view.SetViewportUpdateMode(widgets.QGraphicsView__FullViewportUpdate)
	view.SetVerticalScrollBarPolicy(core.Qt__ScrollBarAlwaysOn)
	view.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)

	view.FitInView(core.NewQRectF4(0, 0, editorViewWidth, 50), core.Qt__IgnoreAspectRatio)

	view.SetResizeAnchor(widgets.QGraphicsView__AnchorUnderMouse)

	// on macOS the gestures don't need to be grabbed because the operating system will detect them for us
	// if we still grab them qt would do the detection instead of macOS
	if runtime.GOOS != "darwin" {
		view.GrabGesture(core.Qt__PinchGesture, 0) // tell qt we want to know about pinch gestures
	}

	view.ConnectWheelEvent(func(event *gui.QWheelEvent) {
		if (event.Modifiers()&core.Qt__AltModifier != 0) && event.Type() == core.QEvent__Wheel {
			delta := event.PixelDelta()
			// TODO: Use AngleDelta on platforms that don't support PixelDelta (See https://doc.qt.io/qt-5/qwheelevent.html#pixelDelta)
			view.Scale(1, 1+float64(delta.Y())/1000)
			return
		}

		// this prevents scrolling on the P Axis
		// we create a new QWheelEvent without the P movement and then call the default wheel event handler with that
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

		event.Ignore() // handle it later
		view.WheelEventDefault(event)
	})

	view.ConnectEvent(func(event *core.QEvent) bool {
		if event.Type() == core.QEvent__Gesture {
			gestureEvent := widgets.NewQGestureEventFromPointer(event.Pointer())
			gesture := gestureEvent.Gesture(core.Qt__PinchGesture)
			if gesture.Pointer() != nil { // not a PinchGesture
				pinchGesture := widgets.NewQPinchGestureFromPointer(gesture.Pointer())
				view.Scale(1, pinchGesture.ScaleFactor())
				return true
			}
		}

		if event.Type() == core.QEvent__NativeGesture {
			gestureEvent := gui.NewQNativeGestureEventFromPointer(event.Pointer())
			if gestureEvent.GestureType() == core.Qt__ZoomNativeGesture {
				view.Scale(1, 1+gestureEvent.Value())
				return true
			}
		}

		event.Ignore()
		return view.EventDefault(event)
	})

	view.ConnectResizeEvent(func(event *gui.QResizeEvent) {
		if event.OldSize().Width() == -1 {
			// ignore if this is the first event at the start
			return
		}

		// calculate the scaling of the window and then apply the same to the view
		view.Scale(
			float64(event.Size().Width())/float64(event.OldSize().Width()),
			float64(event.Size().Height())/float64(event.OldSize().Height()),
		)
	})

	return view
}
