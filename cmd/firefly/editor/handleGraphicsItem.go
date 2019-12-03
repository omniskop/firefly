package editor

import (
	"github.com/omniskop/firefly/pkg/project/vectorpath"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type handleGraphicsItem struct {
	*widgets.QGraphicsPixmapItem
	parent                   *elementGraphicsItem
	index                    int
	ignoreNextPositionChange bool
	moveStartPosition        *core.QPointF
}

func newHandleGraphicsItem(parent *elementGraphicsItem, point vectorpath.Point, index int) *handleGraphicsItem {
	picture := gui.NewQIcon5("assets/images/testPicture/testPicture.png").Pixmap2(25, 25, gui.QIcon__Normal, gui.QIcon__On)

	item := handleGraphicsItem{
		QGraphicsPixmapItem: widgets.NewQGraphicsPixmapItem2(picture, parent),
		parent:              parent,
		index:               index,
		moveStartPosition:   core.NewQPointF(),
	}

	item.SetPos(qtPoint(point.Sub(parent.element.Shape.Origin()))) // convert the scene coordinate to parent coordinates
	item.SetFlags(widgets.QGraphicsItem__ItemIgnoresTransformations | widgets.QGraphicsItem__ItemIsMovable | widgets.QGraphicsItem__ItemSendsScenePositionChanges)

	// connect all necessary events
	item.ConnectItemChange(item.itemChangeEvent)
	item.ConnectMousePressEvent(item.mousePressEvent)

	return &item
}

func (item *handleGraphicsItem) mousePressEvent(event *widgets.QGraphicsSceneMouseEvent) {
	item.moveStartPosition = item.parent.Pos() // save the position of the parent element for later
	// if we would ignore the event here the default implementation would not accept it and thus this element would not react to the event
	//event.Ignore()
	item.MousePressEventDefault(event)
}

func (item *handleGraphicsItem) itemChangeEvent(change widgets.QGraphicsItem__GraphicsItemChange, value *core.QVariant) *core.QVariant {
	if change == widgets.QGraphicsItem__ItemPositionHasChanged {
		if item.ignoreNextPositionChange {
			item.ignoreNextPositionChange = false
			goto end
		}
		// Unfortunately the position of the element will not update correctly because it's position is relative to the
		// parent which could be moving at the same if the handle changes the position of the shape.
		// Pos() is the position relative to where the item got originally picked up.
		// ScenePos() is the current position of the parent plus Pos().
		// Because Pos() does not reflect the real current position we can't use one of the above.
		// That's why we save the position of the parent element from when the handle got picked up
		// and now add Pos(). This will result in the correct absolute position of the handle in scene coordinates.
		trueScenePos := core.NewQPointF3(item.moveStartPosition.X()+item.Pos().X(), item.moveStartPosition.Y()+item.Pos().Y())
		item.parent.element.Shape.SetHandle(item.index, vpPoint(trueScenePos))
		item.parent.updatePath()
		// Due to the whole position thing the default qt item movement does not work anymore.
		// If a handle gets moved down and that also moves the shape (and thus the parent) the movement gets doubled.
		// Once due to qt moving the handle with the mouse and then a second time because the parent moves the handle
		// with it. That's why the element needs to update all the handles including this one.
		// If there isn't a simple fix for this situation it might be best to implement the movement completely on our
		// own instead of partially using the qt implementation and then reverting it's effect.
		// Which is effectively what we do now.
	}

end:
	return item.ItemChangeDefault(change, value)
}

func (item *handleGraphicsItem) updatePosition(point vectorpath.Point) {
	// Setting the position will result in the item chane event to be called.
	// To prevent an infinite loop of handles updating each other we set this to true and check it later in the event.
	item.ignoreNextPositionChange = true
	// The point is in scene coordinates. We need to map the point to parent coordinates
	// This is the same as
	// item.SetPos(item.parent.MapFromScene(qtPoint(point)))
	// but we save a c-call
	item.SetPos(qtPoint(point.Sub(item.parent.element.Shape.Origin())))
}
