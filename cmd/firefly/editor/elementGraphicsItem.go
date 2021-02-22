package editor

import (
	"github.com/sirupsen/logrus"

	"github.com/omniskop/firefly/pkg/project"
	"github.com/omniskop/firefly/pkg/project/vectorpath"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var selectionPen = gui.NewQPen4(
	gui.NewQBrush3(gui.NewQColor3(255, 255, 255, 255), core.Qt__SolidPattern),
	0,
	core.Qt__SolidLine,
	core.Qt__FlatCap,
	core.Qt__BevelJoin,
)

func init() {
	selectionPen.SetCosmetic(true)
}

type elementGraphicsItem struct {
	*widgets.QGraphicsPathItem                       // the underlying path in the QGraphicsScene
	element                    *project.Element      // the element
	parent                     *stage                // the parent editor this element belongs to
	handles                    []*handleGraphicsItem // the handle items that are visible when the element is selected
	gradientItem               *gradientGraphicsItem
	dragStartPosition          *core.QPointF // position of the element when the user started to move it
	ignoreNextPositionChange   byte
}

func newElementGraphicsItem(parentStage *stage, element *project.Element) *elementGraphicsItem {
	item := elementGraphicsItem{
		QGraphicsPathItem: widgets.NewQGraphicsPathItem2(pathFromElement(element), nil),
		element:           element,
		parent:            parentStage,
	}
	item.SetPos(qtPoint(element.Shape.Origin()))
	item.SetZValue(element.ZIndex)
	item.updatePattern()
	item.SetPen(noPen)
	item.SetFlags(widgets.QGraphicsItem__ItemSendsScenePositionChanges | widgets.QGraphicsItem__ItemIsMovable)
	item.ConnectMousePressEvent(item.mousePressEvent)
	item.ConnectItemChange(item.itemChangeEvent)
	return &item
}

func (item *elementGraphicsItem) updatePath() {
	item.ignoreNextPositionChange = 1
	item.PrepareGeometryChange()
	item.SetPos(qtPoint(item.element.Shape.Origin()))
	item.SetPath(pathFromElement(item.element))
	item.SetZValue(item.element.ZIndex)
	item.updateHandles(-1)
}

// updatePattern sets the brush of the element and updates the gradient ui if necessary
func (item *elementGraphicsItem) updatePattern() {
	item.SetBrush(NewQBrushFromPattern(item.element.Pattern)) // TODO: modify brush instead of replacing it
	if item.handles == nil {
		// the element is not selected and we do not need to update a potential gradient
		return
	}
	switch pat := item.element.Pattern.(type) {
	case *project.SolidColor:
		if item.gradientItem != nil {
			item.Scene().RemoveItem(item.gradientItem)
			item.gradientItem = nil
		}
	case *project.LinearGradient:
		if item.gradientItem == nil {
			item.gradientItem = newGradientGraphicsItem(item, pat)
		} else {
			item.gradientItem.updateGradient(pat)
		}
	}
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
	event.Accept() // accept this event to stop this event from propagating to the parent
	item.dragStartPosition = item.Pos()
	item.parent.elementHasBeenClicked(item, event)
	item.MousePressEventDefault(event)
}

func (item *elementGraphicsItem) itemChangeEvent(change widgets.QGraphicsItem__GraphicsItemChange, value *core.QVariant) *core.QVariant {
	if change == widgets.QGraphicsItem__ItemPositionChange {
		if item.parent.creationElement != nil && item.parent.creationElement.Pointer() != item.Pointer() {
			// if an element is currently being created and it is not this element itself the move should be ignored
			// the change will be overwritten by the return value of this function
			return core.NewQVariant28(core.NewQPointF())
		}
		if item.ignoreNextPositionChange == 1 {
			item.ignoreNextPositionChange = 0
			goto end
		}

		newPos := value.ToPointF()

		if item.ignoreNextPositionChange == 2 {
			// This item has been moved because it is part of a selection.
			// We will leave here because we don't need to update other elements.
			item.element.Shape.SetOrigin(vpPoint(newPos))
			item.ignoreNextPositionChange = 0
			goto end
		}

		if gui.QGuiApplication_KeyboardModifiers()&core.Qt__ShiftModifier != 0 {
			// If the user is holding the shift key we will restrict the item movement to only one axis.
			// The axis with the least required change is chosen.
			diffX := newPos.X() - item.dragStartPosition.X()
			diffY := newPos.Y() - item.dragStartPosition.Y()
			pixelDiff := item.parent.mapRelativeFromScene(core.NewQPointF3(diffX, diffY))
			if pixelDiff.X() > pixelDiff.Y() {
				// change on X axis is larger than on Y, we will keep the Y axis constant
				newPos.SetY(item.dragStartPosition.Y())
			} else {
				newPos.SetX(item.dragStartPosition.X())
			}
			value = core.NewQVariant28(newPos) // update value
		}

		item.element.Shape.SetOrigin(vpPoint(newPos))

		// update other selected elements
		oldPos := item.Pos()
		change := core.NewQPointF3(newPos.X()-oldPos.X(), newPos.Y()-oldPos.Y())
		for _, element := range item.parent.selection.elements {
			if element.Pointer() == item.Pointer() {
				continue
			}
			element.ignoreNextPositionChange = 2
			element.MoveBy(change.X(), change.Y())
		}
	}

end:
	return item.ItemChangeDefault(change, value)
}

func (item *elementGraphicsItem) showHandles() {
	if len(item.handles) != 0 {
		logrus.Warn("element already has handles")
		return
	}

	item.handles = make([]*handleGraphicsItem, len(item.element.Shape.Handles()))

	for i, handle := range item.element.Shape.Handles() {
		item.handles[i] = newHandleGraphicsItem(item, handle, i)
	}

	item.SetPen(selectionPen)

	// gradient
	if item.element.Pattern != nil {
		if linearGradient, ok := item.element.Pattern.(*project.LinearGradient); ok {
			item.gradientItem = newGradientGraphicsItem(item, linearGradient)
		}
	}
}

func (item *elementGraphicsItem) hideHandles() {
	scene := item.Scene()
	for _, handleItem := range item.handles {
		scene.RemoveItem(handleItem)
	}
	item.handles = nil

	if item.gradientItem != nil {
		scene.RemoveItem(item.gradientItem)
		item.gradientItem = nil
	}

	item.SetPen(noPen)
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
