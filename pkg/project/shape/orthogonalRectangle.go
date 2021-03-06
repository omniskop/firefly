package shape

import (
	"encoding/json"
	"errors"

	"github.com/omniskop/firefly/pkg/project/vectorpath"
)

// OrthogonalRectangle is a rectangle shape whose edges are orthogonal to the coordinate system
type OrthogonalRectangle struct {
	path vectorpath.Path
}

var _ Shape = (*OrthogonalRectangle)(nil) // make sure OrthogonalRectangle implements the Shape interface

// NewOrthogonalRectangle creates a new shape with the top left position and width and height
func NewOrthogonalRectangle(pos vectorpath.Point, width float64, height float64) *OrthogonalRectangle {
	return &OrthogonalRectangle{
		path: vectorpath.Path{
			Start: pos,
			Segments: []vectorpath.Segment{
				&vectorpath.Line{Point: vectorpath.Point{P: width, T: 0}},
				&vectorpath.Line{Point: vectorpath.Point{P: 0, T: height}},
				&vectorpath.Line{Point: vectorpath.Point{P: -width, T: 0}},
				&vectorpath.Line{Point: vectorpath.Point{P: 0, T: -height}},
			},
		},
	}
}

func NewEmptyOrthogonalRectangle() *OrthogonalRectangle {
	return NewOrthogonalRectangle(vectorpath.Point{}, 0, 0)
}

// Time returns the temporal position of the rectangle in seconds
func (or *OrthogonalRectangle) Time() float64 {
	return or.path.Start.T
}

// Duration returns the duration of the shape in seconds
func (or *OrthogonalRectangle) Duration() float64 {
	return or.path.Duration()
}

// Width returns the width of the shape
func (or *OrthogonalRectangle) Width() float64 {
	return or.path.Segments[0].EndPoint().P
}

// Bounds returns the underlying rectangle
func (or *OrthogonalRectangle) Bounds() vectorpath.Rect {
	return vectorpath.NewRect(
		or.path.Start.P,
		or.path.Start.T,
		or.path.Segments[0].EndPoint().P,
		or.path.Segments[1].EndPoint().T,
	)
}

// Move the rectangle by some amount
func (or *OrthogonalRectangle) Move(by vectorpath.Point) {
	or.path.Start = or.path.Start.Add(by)
}

// Origin returns the top left point of the rectangle
func (or *OrthogonalRectangle) Origin() vectorpath.Point {
	return or.path.Start
}

// SetOrigin sets the top left point of the rectangle to a new point
func (or *OrthogonalRectangle) SetOrigin(l vectorpath.Point) {
	or.path.Start = l
}

// Path returns the path of the shape that should be rendered
func (or *OrthogonalRectangle) Path() vectorpath.Path {
	return or.path
}

// Handles returns all handles for this shape that the user can then use to manipulate the shape
func (or *OrthogonalRectangle) Handles() []vectorpath.Point {
	top := vectorpath.Point{
		P: or.path.Start.P + or.path.Segments[0].EndPoint().P/2,
		T: or.path.Start.T,
	}
	topRight := or.path.Start.Add(or.path.Segments[0].EndPoint())
	right := vectorpath.Point{
		P: topRight.P,
		T: topRight.T + or.path.Segments[1].EndPoint().T/2,
	}
	bottomRight := topRight.Add(or.path.Segments[1].EndPoint())
	bottom := vectorpath.Point{
		P: bottomRight.P + or.path.Segments[2].EndPoint().P/2,
		T: bottomRight.T,
	}
	bottomLeft := bottomRight.Add(or.path.Segments[2].EndPoint())
	left := vectorpath.Point{
		P: bottomLeft.P,
		T: bottomLeft.T + or.path.Segments[3].EndPoint().T/2,
	}
	return []vectorpath.Point{
		or.path.Start,
		topRight,
		bottomRight,
		bottomLeft,
		top,
		right,
		bottom,
		left,
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
		difference := absolutePoint.Sub(or.path.Start)
		or.path.Start = absolutePoint

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
		absoluteHandlePos := or.path.PointAfter(1)

		// difference between the handle positions
		difference := absolutePoint.Sub(absoluteHandlePos)

		or.path.Start.T = absolutePoint.T
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
		absoluteHandlePos := or.path.PointAfter(2)

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
		absoluteHandlePos := or.path.PointAfter(3)

		// difference between the handle positions
		difference := absolutePoint.Sub(absoluteHandlePos)

		or.path.Start.P = absolutePoint.P
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
	case 4: // move the top handle
		difference := absolutePoint.Sub(or.path.Start)
		or.path.Start.T = absolutePoint.T

		or.path.Segments[1].Move(vectorpath.Point{ // right
			P: 0,
			T: -difference.T,
		})

		or.path.Segments[3].Move(vectorpath.Point{ // left
			P: 0,
			T: difference.T,
		})
	case 5: // move the right handle
		// absolute position of the original handle
		absoluteHandlePos := or.path.PointAfter(1)

		// difference between the handle positions
		difference := absolutePoint.Sub(absoluteHandlePos)

		or.path.Segments[0].Move(vectorpath.Point{
			T: 0,
			P: difference.P,
		}) // top right

		or.path.Segments[2].Move(vectorpath.Point{ // bottom left
			T: 0,
			P: -difference.P,
		})
	case 6: // move the bottom handle
		// absolute position of the original handle
		absoluteHandlePos := or.path.PointAfter(2)

		// difference between the handle positions
		difference := absolutePoint.Sub(absoluteHandlePos)

		or.path.Segments[1].Move(vectorpath.Point{ // bottom right
			T: difference.T,
			P: 0,
		})

		or.path.Segments[3].Move(vectorpath.Point{ // top left
			T: -difference.T,
			P: 0,
		})
	case 7: // move the left handle
		// absolute position of the original handle
		absoluteHandlePos := or.path.PointAfter(3)

		// difference between the handle positions
		difference := absolutePoint.Sub(absoluteHandlePos)

		or.path.Start.P = absolutePoint.P
		or.path.Segments[0].Move(vectorpath.Point{ // top right
			T: 0,
			P: -difference.P,
		})

		or.path.Segments[2].Move(vectorpath.Point{ // bottom left
			T: 0,
			P: difference.P,
		})
	}
}

func (or *OrthogonalRectangle) SetCreationBounds(origin vectorpath.Point, size vectorpath.Point) {
	if size.P < 0 {
		origin.P += size.P
		size.P = -size.P
	}
	if size.T < 0 {
		origin.T += size.T
		size.T = -size.T
	}
	or.path.Start = origin
	or.SetHandle(2, origin.Add(size))
}

func (or *OrthogonalRectangle) MirrorP() {
	// No actions required
}

func (or *OrthogonalRectangle) Copy() Shape {
	return NewOrthogonalRectangle(or.Origin(), or.Width(), or.Duration())
}

func (or *OrthogonalRectangle) MarshalJSON() ([]byte, error) {
	var values = map[string]interface{}{
		"__TYPE__": "OrthogonalRectangle",
		"Shape": map[string]vectorpath.Point{
			"Position":   or.path.Start,
			"Dimensions": or.path.PointAfter(2).Sub(or.path.Start),
		},
	}
	return json.Marshal(values)
}

func (or *OrthogonalRectangle) UnmarshalJSON(raw []byte) error {
	var values = make(map[string]vectorpath.Point)
	err := json.Unmarshal(raw, &values)
	if err != nil {
		return err
	}

	pos, ok := values["Position"]
	if !ok {
		return errors.New("orthogonal rectangle has missing key 'Position'")
	}
	dim, ok := values["Dimensions"]
	if !ok {
		return errors.New("orthogonal rectangle has missing key 'Dimensions'")
	}

	or.path = vectorpath.Path{
		Start: pos,
		Segments: []vectorpath.Segment{
			&vectorpath.Line{Point: vectorpath.Point{P: dim.P, T: 0}},
			&vectorpath.Line{Point: vectorpath.Point{P: 0, T: dim.T}},
			&vectorpath.Line{Point: vectorpath.Point{P: -dim.P, T: 0}},
			&vectorpath.Line{Point: vectorpath.Point{P: 0, T: -dim.T}},
		},
	}
	return nil
}
