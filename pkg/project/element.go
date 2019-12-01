package project

import "github.com/omniskop/firefly/pkg/project/vectorpath"

// An Element that is located at a specific point in time in the scene
type Element struct {
	ZIndex  float64 // a coordinate relative to other elements in the scene. Higher numbers will be drawn on top of lower ones
	Shape   Shape   // the actual visual shape of the element
	Pattern Pattern // the pattern that fills the body of the element
}

// MapLocalToRelative maps local coordinates to relative coordinates.
//
// Local coordinates are at (0,0) at the top left of the bounding box of the shape and (1,1) at the bottom right.
// Relative coordinates have their origin at the top left of the bounding box of the shape but their scale is the same as the scene.
func (e *Element) MapLocalToRelative(local vectorpath.Point) vectorpath.Point {
	return vectorpath.Point{
		P: e.Shape.Width() * local.P,
		T: e.Shape.Duration() * local.T,
	}
}

// MapRelativeToLocal maps relative coordinates to local coordinates.
// See MapLocalToRelative for an explanation of the coordinates.
func (e *Element) MapRelativeToLocal(relative vectorpath.Point) vectorpath.Point {
	// TODO: maybe prevent width and duration from every being zero?
	width := e.Shape.Width()
	duration := e.Shape.Duration()
	if width == 0 || duration == 0 {
		return vectorpath.Point{P: 0, T: 0}
	}
	return vectorpath.Point{
		P: relative.P / width,
		T: relative.T / duration,
	}
}
