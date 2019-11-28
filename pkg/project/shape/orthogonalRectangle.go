package shape

import (
	"github.com/omniskop/firefly/pkg/project"
	"github.com/omniskop/firefly/pkg/project/vectorpath"
)

// OrthogonalRectangle is a rectangle shape whose edges are orthogonal to the coordinate system
type OrthogonalRectangle struct {
	t    float64
	path vectorpath.Path
}

var _ project.Shape = (*OrthogonalRectangle)(nil) // make sure OrthogonalRectangle implements the Shape interface

// NewOrthogonalRectangle creates a new shape with the top left position and width and height
func NewOrthogonalRectangle(pos vectorpath.Point, width float64, height float64) *OrthogonalRectangle {
	return &OrthogonalRectangle{
		t: pos.T,
		path: vectorpath.Path{
			P: pos.P,
			Segments: []vectorpath.Segment{
				&vectorpath.Line{Point: vectorpath.Point{P: width, T: 0}},
				&vectorpath.Line{Point: vectorpath.Point{P: 0, T: height}},
				&vectorpath.Line{Point: vectorpath.Point{P: -width, T: 0}},
				&vectorpath.Line{Point: vectorpath.Point{P: 0, T: -height}},
			},
		},
	}
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
	if absolutePoint.P < 0 {
		absolutePoint.P = 0
	} else if absolutePoint.P > 1 {
		absolutePoint.P = 1
	}
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
		absoluteHandlePos := or.path.PointAfter(or.t, 3)

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

// Move the rectangle by some amount
func (or *OrthogonalRectangle) Move(by vectorpath.Point) {
	or.t += by.T
	or.path.P += by.P
}

// Location returns the absolute location of the shape
func (or *OrthogonalRectangle) Location() vectorpath.Point {
	return vectorpath.Point{P: or.path.P, T: or.t}
}

// SetLocation will set the location of the shape to a new position
func (or *OrthogonalRectangle) SetLocation(point vectorpath.Point) {
	or.t = point.T
	or.path.P = point.P
}
