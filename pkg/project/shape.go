package project

import (
	"github.com/omniskop/firefly/pkg/project/vectorpath"
)

// Shape is the visual part of an element on the scene
type Shape interface {
	Time() float64
	Duration() float64
	Path() vectorpath.Path
	Handles() []vectorpath.Point
	SetHandle(int, vectorpath.Point)
	Move(vectorpath.Point)
	Location() vectorpath.Point
	SetLocation(vectorpath.Point)
}
