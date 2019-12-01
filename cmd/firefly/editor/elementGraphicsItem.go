package editor

import (
	"github.com/sirupsen/logrus"

	"github.com/omniskop/firefly/pkg/project"
	"github.com/omniskop/firefly/pkg/project/vectorpath"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type elementGraphicsItem struct {
	*widgets.QGraphicsPathItem                       // the underlying path in the QGraphicsScene
	element                    *project.Element      // the element
	editor                     *Editor               // the parent editor this element belongs to
	handles                    []*handleGraphicsItem // the handle items that are visible when the element is selected
	gradientItem               *gradientGraphicsItem
	ignoreNextPositionChange   bool
}

func newElementGraphicsItem(editor *Editor, element *project.Element) *elementGraphicsItem {
	item := elementGraphicsItem{
		QGraphicsPathItem: widgets.NewQGraphicsPathItem2(pathFromElement(element), nil),
		element:           element,
		editor:            editor,
	}
	item.SetPos(qtPoint(element.Shape.Location()))
	item.updatePattern()
	item.SetPen(noPen)
	item.SetFlags(widgets.QGraphicsItem__ItemSendsScenePositionChanges | widgets.QGraphicsItem__ItemIsMovable)
	item.ConnectMousePressEvent(item.mousePressEvent)
	item.ConnectItemChange(item.itemChangeEvent)
	return &item
}

func (item *elementGraphicsItem) updatePath() {
	item.ignoreNextPositionChange = true
	item.SetPos(qtPoint(item.element.Shape.Location()))
	item.SetPath(pathFromElement(item.element))
	item.updateHandles(-1)
}

// updatePattern sets the brush of the element
func (item *elementGraphicsItem) updatePattern() {
	item.SetBrush(NewQBrushFromPattern(item.element.Pattern)) // TODO: modify brush instead of replacing it
}

func (item *elementGraphicsItem) updateHandles(except int) {
	shapeHandles := item.element.Shape.Handles()
	for i, handle := range item.handles {
		if i == except {
			continue
		}
		handle.updatePosition(shapeHandles[i])
	}

	// update the gradient graphics item
	// i might move this into it's own method at some point
	if item.gradientItem != nil {
		item.gradientItem.updateShape(-100)
	}
}

func (item *elementGraphicsItem) mousePressEvent(event *widgets.QGraphicsSceneMouseEvent) {
	logrus.Trace("element pressed")
	item.selectElement()
}

func (item *elementGraphicsItem) itemChangeEvent(change widgets.QGraphicsItem__GraphicsItemChange, value *core.QVariant) *core.QVariant {
	if change == widgets.QGraphicsItem__ItemPositionHasChanged {
		if item.ignoreNextPositionChange {
			item.ignoreNextPositionChange = false
			goto end
		}
		item.element.Shape.SetLocation(vpPoint(item.ScenePos()))
	}
	//if change == widgets.QGraphicsItem__ItemPositionChange {
	//	newPos := core.NewQPointFFromPointer(value.Pointer())
	//	diff := vpPoint(newPos).Sub(vpPoint(item.Pos()))
	//	fmt.Println(diff)
	//	item.element.Shape.Move(diff)
	//}

end:
	return item.ItemChangeDefault(change, value)
}

func (item *elementGraphicsItem) selectElement() {
	if len(item.handles) != 0 {
		logrus.Trace("element already selected")
		return
	}

	item.handles = make([]*handleGraphicsItem, len(item.element.Shape.Handles()))

	for i, handle := range item.element.Shape.Handles() {
		item.handles[i] = newHandleGraphicsItem(item, handle, i)
	}
	logrus.WithField("handles", len(item.handles)).Trace("created handles")

	// gradient
	if item.element.Pattern != nil {
		if linearGradient, ok := item.element.Pattern.(*project.LinearGradient); ok {
			item.gradientItem = newGradientGraphicsItem(item, linearGradient)
		}
	}

	item.editor.elementSelected(item)
}

func (item *elementGraphicsItem) deselectElement() {
	logrus.WithFields(logrus.Fields{"handles": len(item.handles), "item": item}).Trace("deselectElement called")
	scene := item.Scene()
	for _, handleItem := range item.handles {
		scene.RemoveItem(handleItem)
	}
	item.handles = nil

	if item.gradientItem != nil {
		scene.RemoveItem(item.gradientItem)
		item.gradientItem = nil
	}
}

func pathFromElement(element *project.Element) *gui.QPainterPath {
	elementPath := element.Shape.Path()
	path := gui.NewQPainterPath()
	currentPosition := vectorpath.Point{P: 0, T: 0}
	path.MoveTo(qtPoint(currentPosition))

	for _, segment := range elementPath.Segments {
		currentPosition = addSegmentToPath(path, currentPosition, segment)
	}

	return path
}

func addSegmentToPath(path *gui.QPainterPath, startingPoint vectorpath.Point, segment vectorpath.Segment) vectorpath.Point {
	switch seg := segment.(type) {
	case *vectorpath.Line:
		path.LineTo(qtPoint(seg.Point.Add(startingPoint)))
	case *vectorpath.CubicCurve:
		path.CubicTo(qtPoint(seg.ControlA.Add(startingPoint)), qtPoint(seg.ControlB.Add(startingPoint)), qtPoint(seg.End.Add(startingPoint)))
	case *vectorpath.QuadCurve:
		path.QuadTo(qtPoint(seg.Control.Add(startingPoint)), qtPoint(seg.End.Add(startingPoint)))
	}
	return startingPoint.Add(segment.EndPoint())
}

func qtPointOffset(p vectorpath.Point, offsetP float64, offsetT float64) *core.QPointF {
	if verticalTimeAxis {
		return core.NewQPointF3(p.P+offsetP*editorViewWidth, p.T+offsetT)
	}
	return core.NewQPointF3(p.T+offsetT, p.P+offsetP*editorViewWidth)
}

func qtPoint(p vectorpath.Point) *core.QPointF {
	if verticalTimeAxis {
		return core.NewQPointF3(p.P*editorViewWidth, p.T)
	}
	return core.NewQPointF3(p.T, p.P*editorViewWidth)
}

func vpPoint(p *core.QPointF) vectorpath.Point {
	if verticalTimeAxis {
		return vectorpath.Point{P: p.X() / editorViewWidth, T: p.Y()}
	}
	return vectorpath.Point{P: p.Y(), T: p.X() / editorViewWidth}
}
