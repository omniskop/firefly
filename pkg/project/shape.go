package project

import "github.com/omniskop/firefly/pkg/project/vectorpath"

// Shape is the visual part of an element on the scene
type Shape interface {
	Time() float64
	Duration() float64
	Path() vectorpath.Path
	Handles() []vectorpath.Point
	SetHandle(int, vectorpath.Point)
}

// OrthogonalRectangle is a rectangle shape whose edges are orthogonal to the coordinate system
type OrthogonalRectangle struct {
	t    float64
	path vectorpath.Path
}

// Time returns the temporal position of the rectangle in seconds
func (or *OrthogonalRectangle) Time() float64 {
	return or.t
}

// Duration returns the duration of the shape in seconds
func (or *OrthogonalRectangle) Duration() float64 {
	return or.path.Duration()
}

// Path returns the path of the shape that should be rendered
func (or *OrthogonalRectangle) Path() vectorpath.Path {
	return or.path
}

// Handles returns all handles for this shape that the user can then use to manipulate the shape
func (or *OrthogonalRectangle) Handles() []vectorpath.Point {
	topLeft := vectorpath.Point{
		T: or.t,
		P: or.path.P,
	}
	topRight := topLeft.Add(or.path.Segments[0].EndPoint())
	bottomRight := topRight.Add(or.path.Segments[1].EndPoint())
	bottomLeft := bottomRight.Add(or.path.Segments[2].EndPoint())
	return []vectorpath.Point{
		topLeft,
		topRight,
		bottomRight,
		bottomLeft,
	}
}

// SetHandle receives the index of the handle that should be changed and the new point value
func (or *OrthogonalRectangle) SetHandle(i int, absolutePoint vectorpath.Point) {
	switch i {
	case 0: // move the top left handle
		difference := absolutePoint.Sub(vectorpath.Point{
			T: or.t,
			P: or.path.P,
		})
		or.t = absolutePoint.T
		or.path.P = absolutePoint.P

		or.path.Segments[0].Move(vectorpath.Point{ // top right
			P: -difference.P, // counteract the movement of the start position
			T: 0,             // does not move relative to the start on the time axis
		})

		or.path.Segments[1].Move(vectorpath.Point{ // bottom right
			P: 0,
			T: -difference.T,
		})

		or.path.Segments[2].Move(vectorpath.Point{ // bottom left
			P: difference.P,
			T: 0,
		})

		or.path.Segments[3].Move(vectorpath.Point{ // top left
			P: 0,
			T: difference.T,
		})
	case 1: // move the top right handle
		// absolute position of the original handle
		absoluteHandlePos := or.path.PointAfter(or.t, 1)

		// difference between the handle positions
		difference := absolutePoint.Sub(absoluteHandlePos)

		or.t = absolutePoint.T
		or.path.Segments[0].Move(vectorpath.Point{
			T: 0,
			P: difference.P,
		}) // top right

		or.path.Segments[1].Move(vectorpath.Point{ // bottom right
			T: -difference.T,
			P: 0,
		})

		or.path.Segments[2].Move(vectorpath.Point{ // bottom left
			T: 0,
			P: -difference.P,
		})

		or.path.Segments[3].Move(vectorpath.Point{ // top left
			T: difference.T,
			P: 0,
		})
	case 2: // move the bottom right handle
		// absolute position of the original handle
		absoluteHandlePos := or.path.PointAfter(or.t, 2)

		// difference between the handle positions
		difference := absolutePoint.Sub(absoluteHandlePos)

		or.path.Segments[0].Move(vectorpath.Point{ // top right
			T: 0,
			P: difference.P,
		})

		or.path.Segments[1].Move(vectorpath.Point{ // bottom right
			T: difference.T,
			P: 0,
		})

		or.path.Segments[2].Move(vectorpath.Point{ // bottom left
			T: 0,
			P: -difference.P,
		})

		or.path.Segments[3].Move(vectorpath.Point{ // top left
			T: -difference.T,
			P: 0,
		})
	case 3: // move the bottom left handle
		// absolute position of the original handle
		absoluteHandlePos := or.path.PointAfter(or.t, 2)

		// difference between the handle positions
		difference := absolutePoint.Sub(absoluteHandlePos)

		or.path.P = absolutePoint.P
		or.path.Segments[0].Move(vectorpath.Point{ // top right
			T: 0,
			P: -difference.P,
		})

		or.path.Segments[1].Move(vectorpath.Point{ // bottom right
			T: difference.T,
			P: 0,
		})

		or.path.Segments[2].Move(vectorpath.Point{ // bottom left
			T: 0,
			P: difference.P,
		})

		or.path.Segments[3].Move(vectorpath.Point{ // top left
			T: -difference.T,
			P: 0,
		})
	}
}
