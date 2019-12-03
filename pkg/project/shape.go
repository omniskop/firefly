package project

import (
	"github.com/omniskop/firefly/pkg/project/vectorpath"
)

// Shape is the visual part of an element on the scene
type Shape interface {
	Time() float64              // point in time when the shape start
	Duration() float64          // duration of time that the shape takes up
	Width() float64             // visual width of the shape
	Bounds() vectorpath.Rect    // outer rectangular bounds of the shape
	Move(vectorpath.Point)      // move the shape by some amount
	Origin() vectorpath.Point   // get the point where the path of the shape starts (does not have to be the same as Bounds().Location)
	SetOrigin(vectorpath.Point) // set the origin of the shape

	Path() vectorpath.Path
	Handles() []vectorpath.Point     // returns all points where the user can manipulate the shape
	SetHandle(int, vectorpath.Point) // set new position of a handle
}
