package editor

import (
	"image/color"

	"github.com/sirupsen/logrus"

	"github.com/omniskop/firefly/pkg/project"
	"github.com/omniskop/firefly/pkg/project/vectorpath"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var gradientLinePen = gui.NewQPen4(
	gui.NewQBrush3(gui.NewQColor3(255, 255, 255, 255), core.Qt__SolidPattern),
	5,
	core.Qt__SolidLine,
	core.Qt__FlatCap,
	core.Qt__BevelJoin,
)

func init() {
	gradientLinePen.SetCosmetic(true) // prevent the pen from being transformed
}

type gradientGraphicsItem struct {
	*widgets.QGraphicsLineItem
	parent   *elementGraphicsItem
	gradient *project.LinearGradient
	start    *colorDotItem
	stop     *colorDotItem
}

func newGradientGraphicsItem(parent *elementGraphicsItem, gradient *project.LinearGradient) *gradientGraphicsItem {
	line := core.NewQLineF2(
		qtPoint(parent.element.MapLocalToRelative(gradient.Start.Point)),
		qtPoint(parent.element.MapLocalToRelative(gradient.Stop.Point)),
	)
	item := gradientGraphicsItem{
		QGraphicsLineItem: widgets.NewQGraphicsLineItem2(line, parent),
		parent:            parent,
		gradient:          gradient,
	}

	item.start = newColorDotGraphicsItem(&item, parent.element.MapLocalToRelative(gradient.Start.Point), gradient.Start.Color, -1)
	item.stop = newColorDotGraphicsItem(&item, parent.element.MapLocalToRelative(gradient.Stop.Point), gradient.Stop.Color, -2)

	item.SetPen(gradientLinePen)

	// enabling this shadow has a significant performance impact on objects that take ob a lot of screen space
	/*shadow := widgets.NewQGraphicsDropShadowEffect(nil)
	shadow.SetColor(gui.NewQColor3(0, 0, 0, 165))
	shadow.SetBlurRadius(6)
	shadow.SetOffset2(0, 1)
	item.SetGraphicsEffect(shadow)*/

	return &item
}

func (item *gradientGraphicsItem) updateShape(except int) {
	// TODO: change line instead of replacing it
	item.QGraphicsLineItem.SetLine(core.NewQLineF2(
		qtPoint(item.parent.element.MapLocalToRelative(item.gradient.Start.Point)),
		qtPoint(item.parent.element.MapLocalToRelative(item.gradient.Stop.Point)),
	))
	item.start.update(item.parent.element.MapLocalToRelative(item.gradient.Start.Point), item.gradient.Start.Color, -1)
	item.stop.update(item.parent.element.MapLocalToRelative(item.gradient.Stop.Point), item.gradient.Stop.Color, -2)
}

func (item *gradientGraphicsItem) updateGradient(gradient *project.LinearGradient) {
	item.gradient = gradient
	item.start.update(item.parent.element.MapLocalToRelative(item.gradient.Start.Point), gradient.Start.Color, -1)
	item.stop.update(item.parent.element.MapLocalToRelative(item.gradient.Stop.Point), gradient.Stop.Color, -2)
}

func (item *gradientGraphicsItem) mousePressEvent(event *widgets.QGraphicsSceneMouseEvent) {
	// if we would ignore the event here the default implementation would not accept it and thus this element would not react to the event
	//event.Ignore()
	item.MousePressEventDefault(event)
}

func (item *gradientGraphicsItem) itemChangeEvent(change widgets.QGraphicsItem__GraphicsItemChange, value *core.QVariant) *core.QVariant {
	if change == widgets.QGraphicsItem__ItemPositionHasChanged {

	}

	return item.ItemChangeDefault(change, value)
}

func (item *gradientGraphicsItem) SetStopPosition(index int, point vectorpath.Point) {
	point = item.parent.element.MapRelativeToLocal(point)
	if index == -1 { // Start
		item.gradient.Start.Point = point
	} else if index == -2 { // Stop
		item.gradient.Stop.Point = point
	} else if index < len(item.gradient.Steps) { // color step
		item.gradient.Steps[index].Position = (point.P - item.gradient.Start.P) / (item.gradient.Stop.P - item.gradient.Start.P)
	} else {
		logrus.
			WithFields(logrus.Fields{"index": index, "point": point}).
			Error("gradientGraphicsItem tried to set the position of an unknown color step")
	}
	item.updateShape(index)
	item.parent.updatePattern()
}

type colorDotItem struct {
	*widgets.QGraphicsEllipseItem
	parent                   *gradientGraphicsItem
	index                    int
	ignoreNextPositionChange bool
}

// newColorDotGraphicsItem creates a new colorDotItem
// from the parent gradient, the position of the dot in local coordinates,
// the color of the dot and the index of the dot in the gradient.
func newColorDotGraphicsItem(parent *gradientGraphicsItem, pos vectorpath.Point, col color.Color, index int) *colorDotItem {
	const diameter float64 = 15

	item := colorDotItem{
		QGraphicsEllipseItem: widgets.NewQGraphicsEllipseItem2(
			core.NewQRectF2(
				core.NewQPointF3(-diameter/2, -diameter/2),
				core.NewQSizeF3(diameter, diameter),
			),
			parent,
		),
		parent: parent,
		index:  index,
	}

	item.SetFlags(widgets.QGraphicsItem__ItemIgnoresTransformations | widgets.QGraphicsItem__ItemIsMovable | widgets.QGraphicsItem__ItemSendsScenePositionChanges)
	item.SetPos(qtPoint(pos))

	// set pen
	item.SetPen(gradientLinePen)

	// set brush
	item.SetBrush(gui.NewQBrush3(NewQColorFromColor(col), core.Qt__SolidPattern))

	// signals
	item.ConnectItemChange(item.itemChangeEvent)

	return &item
}

func (item *colorDotItem) update(pos vectorpath.Point, col color.Color, index int) {
	newPos := qtPoint(pos)
	if item.X() != newPos.X() || item.Y() != newPos.Y() {
		item.ignoreNextPositionChange = true
		item.SetPos(newPos)
	}
	brush := item.Brush()
	brush.SetColor(NewQColorFromColor(col))
	item.SetBrush(brush)
	item.index = index
}

func (item *colorDotItem) itemChangeEvent(change widgets.QGraphicsItem__GraphicsItemChange, value *core.QVariant) *core.QVariant {
	if change == widgets.QGraphicsItem__ItemPositionHasChanged {
		if item.ignoreNextPositionChange {
			item.ignoreNextPositionChange = false
			goto end
		}
		item.parent.SetStopPosition(item.index, vpPoint(item.Pos()))
	}

end:
	return item.ItemChangeDefault(change, value)
}
