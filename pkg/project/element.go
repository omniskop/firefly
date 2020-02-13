package project

import (
	"encoding/json"

	"github.com/omniskop/firefly/pkg/project/shape"
	"github.com/omniskop/firefly/pkg/project/vectorpath"
)

// An Element that is located at a specific point in time in the scene
type Element struct {
	ZIndex  float64     // a coordinate relative to other elements in the scene. Higher numbers will be drawn on top of lower ones
	Shape   shape.Shape // the actual visual shape of the element
	Pattern Pattern     // the pattern that fills the body of the element
}

// MapLocalToRelative maps local coordinates to relative coordinates.
//
// Local coordinates are at (0,0) at the top left of the bounding box of the shape and (1,1) at the bottom right.
// Relative coordinates have their origin at the top left of the bounding box of the shape but their scale is the same as the scene.
func (e *Element) MapLocalToRelative(local vectorpath.Point) vectorpath.Point {
	bounds := e.Shape.Bounds()
	return vectorpath.Point{
		P: bounds.Dimensions.P * local.P,
		T: bounds.Dimensions.T * local.T,
	}.Add(bounds.Location.Sub(e.Shape.Origin())) // in the case that the bounds location and the origin of the shape don't align
}

// MapRelativeToLocal maps relative coordinates to local coordinates.
// See MapLocalToRelative for an explanation of the coordinates.
func (e *Element) MapRelativeToLocal(relative vectorpath.Point) vectorpath.Point {
	// TODO: maybe prevent width and duration from every being zero?
	bounds := e.Shape.Bounds()
	if bounds.Dimensions.P == 0 || bounds.Dimensions.T == 0 {
		return vectorpath.Point{P: 0, T: 0}
	}
	relative = relative.Add(e.Shape.Origin().Sub(bounds.Location)) // in the case that the bounds location and the origin of the shape don't align
	return vectorpath.Point{
		P: relative.P / bounds.Dimensions.P,
		T: relative.T / bounds.Dimensions.T,
	}
}

// Copy returns a deep copy of the element
func (e *Element) Copy() *Element {
	return &Element{
		ZIndex:  e.ZIndex,
		Shape:   e.Shape.Copy(),
		Pattern: e.Pattern.Copy(),
	}
}

// MirrorP mirrors this element on the P axis around the center of the scene
func (e *Element) Mirror() {
	e.Shape.MirrorP()
	e.Pattern.MirrorP()

	bounds := e.Shape.Bounds()
	topRight := bounds.Location.P + bounds.Dimensions.P
	mirrored := 0.5 + (0.5 - topRight)
	e.Shape.Move(vectorpath.Point{P: mirrored - bounds.Location.P, T: 0})
}

// UnmarshalJSON will take data and try to parse it into an Element.
// It takes care of handling the Shape and Pattern interfaces with their respective Unmarshal functions.
func (e *Element) UnmarshalJSON(data []byte) error {
	values := make(map[string]*json.RawMessage)
	err := json.Unmarshal(data, &values)
	if err != nil {
		return err
	}

	shapeValue, err := shape.Unmarshal(*values["Shape"])
	if err != nil {
		return err
	}
	e.Shape = shapeValue

	pattern, err := UnmarshalPattern(*values["Pattern"])
	if err != nil {
		return err
	}
	e.Pattern = pattern

	err = json.Unmarshal(*values["ZIndex"], &e.ZIndex)
	return err
}
