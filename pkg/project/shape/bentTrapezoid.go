package shape

import (
	"encoding/json"
	"errors"
	"math"

	"github.com/omniskop/firefly/pkg/project/vectorpath"
)

// BentTrapezoid is a Trapezoid with a curve. The two parallel lines are orthogonal to the time axis and the other two lines are bent.
/*
   position    topWidth
            +-------------+          -
           /               \         |
          /                 \        | duration
         /                   \       |
        /                     \      |
       +-----------------------+     -
              bottomWidth
       |----|
       bottomOffset

*/
type BentTrapezoid struct {
	position     vectorpath.Point
	topWidth     float64
	bottomWidth  float64
	bottomOffset float64
	duration     float64
	bend         vectorpath.Point
}

var _ Shape = (*BentTrapezoid)(nil) // make sure BentTrapezoid implements the Shape interface

// NewBentTrapezoid returns a new BentTrapezoid
func NewBentTrapezoid(topPosition vectorpath.Point, bottomPosition vectorpath.Point, topWidth float64, bottomWidth float64) *BentTrapezoid {
	return &BentTrapezoid{
		position:     topPosition,
		topWidth:     topWidth,
		bottomWidth:  bottomWidth,
		bottomOffset: bottomPosition.P - topPosition.P,
		duration:     bottomPosition.T - topPosition.T,
		bend:         vectorpath.Point{P: 0.5, T: .5},
	}
}

func NewEmptyBentTrapezoid() *BentTrapezoid {
	return NewBentTrapezoid(vectorpath.Point{}, vectorpath.Point{}, 0, 0)
}

// Time returns the point in time where the shape start
func (b *BentTrapezoid) Time() float64 {
	return b.position.T
}

// Duration returns the duration that the shape takes up
func (b *BentTrapezoid) Duration() float64 {
	return b.duration
}

// Width returns the visual width of the shape
func (b *BentTrapezoid) Width() float64 {
	if b.bottomOffset >= 0 {
		return math.Max(b.topWidth, b.bottomOffset+b.bottomWidth)
	}
	return math.Max(-b.bottomOffset+b.topWidth, b.bottomWidth)
}

// Bounds returns the outer bounds of the shape
func (b *BentTrapezoid) Bounds() vectorpath.Rect {
	if b.bottomOffset >= 0 {
		return vectorpath.NewRect(
			b.position.P,
			b.position.T,
			math.Max(b.topWidth, b.bottomOffset+b.bottomWidth),
			b.duration,
		)
	}
	return vectorpath.NewRect(
		b.position.P+b.bottomOffset,
		b.position.T,
		math.Max(-b.bottomOffset+b.topWidth, b.bottomWidth),
		b.duration,
	)
}

// Move the shape by some amount
func (b *BentTrapezoid) Move(offset vectorpath.Point) {
	b.position = b.position.Add(offset)
}

// Origin returns the top left point of the trapezoid
func (b *BentTrapezoid) Origin() vectorpath.Point {
	return b.position
}

// SetOrigin sets the position of the top left point
func (b *BentTrapezoid) SetOrigin(point vectorpath.Point) {
	b.position = point
}

// Path returns a path that describes the bent trapezoid
func (b *BentTrapezoid) Path() vectorpath.Path {
	return vectorpath.Path{
		Start: b.position,
		Segments: []vectorpath.Segment{
			&vectorpath.Line{Point: vectorpath.Point{P: b.topWidth, T: 0}}, // top right
			//&vectorpath.Line{Point: vectorpath.Point{ // bottom right
			//	P: (b.bottomOffset + b.bottomWidth) - (b.topWidth),
			//	T: b.duration,
			//}},
			&vectorpath.QuadCurve{ // bottom right
				Control: vectorpath.Point{
					P: b.interpolate((b.bottomOffset+b.bottomWidth)-b.topWidth, b.bend.P),
					T: b.duration * b.bend.T,
				},
				End: vectorpath.Point{
					P: (b.bottomOffset + b.bottomWidth) - (b.topWidth),
					T: b.duration,
				},
			},
			&vectorpath.Line{Point: vectorpath.Point{ // bottom left
				P: -b.bottomWidth,
				T: 0,
			}},
			//&vectorpath.Line{Point: vectorpath.Point{ // top left
			//	P: -b.bottomOffset,
			//	T: -b.duration,
			//}},
			&vectorpath.QuadCurve{ // bottom right
				Control: vectorpath.Point{
					P: -b.interpolate(b.bottomOffset, 1-b.bend.P),
					T: -b.duration * (1 - b.bend.T),
				},
				End: vectorpath.Point{ // top left
					P: -b.bottomOffset,
					T: -b.duration,
				},
			},
		},
	}
}

func (b *BentTrapezoid) interpolate(a, p float64) float64 {
	if a > 0 {
		return a * p
	}
	return a * (1 - p)
}

// Handles returns the position of all handles
func (b *BentTrapezoid) Handles() []vectorpath.Point {
	return []vectorpath.Point{
		{ // top center
			P: b.position.P + b.topWidth/2,
			T: b.position.T,
		},
		{ // bottom center
			P: b.position.P + b.bottomOffset + b.bottomWidth/2,
			T: b.position.T + b.duration,
		},
		{ // top right
			P: b.position.P + b.topWidth,
			T: b.position.T,
		},
		{ // bottom right
			P: b.position.P + b.bottomOffset + b.bottomWidth,
			T: b.position.T + b.duration,
		},
		{ // bend
			P: b.position.P +
				b.topWidth*b.bend.P +
				b.bottomOffset*b.bend.T +
				(b.bottomWidth-b.topWidth)*b.bend.P*b.bend.T,
			T: b.position.T + b.duration*b.bend.T,
		},
	}
}

// SetHandle sets the handle 'index' to the new point
func (b *BentTrapezoid) SetHandle(index int, absolutePoint vectorpath.Point) {
	switch index {
	case 0: // top - top center
		absolutePoint.T = math.Min(absolutePoint.T, b.position.T+b.duration) // prevent T from beeing after the end of this shape
		diff := absolutePoint.Sub(vectorpath.Point{
			P: b.position.P + b.topWidth/2,
			T: b.position.T,
		})
		b.position = b.position.Add(diff)
		b.bottomOffset -= diff.P
		b.duration -= diff.T
	case 1: // bottom - bottom center
		diff := absolutePoint.Sub(vectorpath.Point{
			P: b.position.P + b.bottomOffset + b.bottomWidth/2,
			T: b.position.T + b.duration,
		})
		b.bottomOffset += diff.P
		b.duration += diff.T
	case 2: // top width - top right
		pDiff := absolutePoint.P - (b.position.P + b.topWidth)
		if b.topWidth+pDiff < 0 {
			pDiff = -b.topWidth
		}
		b.topWidth += pDiff
		b.position.P -= pDiff / 2
		b.bottomOffset += pDiff / 2
	case 3: // bottom width - bottom right
		pDiff := absolutePoint.P - (b.position.P + b.bottomOffset + b.bottomWidth)
		if b.bottomWidth+pDiff < 0 {
			pDiff = -b.bottomWidth
		}
		b.bottomWidth += pDiff
		b.bottomOffset -= pDiff / 2
	case 4: // bend - center
		b.bend.T = (absolutePoint.T - b.position.T) / ((b.position.T + b.duration) - b.position.T)
		b.bend.T = math.Min(1, math.Max(0, b.bend.T))
		b.bend.P = (absolutePoint.P - b.position.P - b.bottomOffset*b.bend.T) / (b.topWidth + (b.bottomWidth-b.topWidth)*b.bend.T)
		b.bend.P = math.Min(1, math.Max(0, b.bend.P))
	}
}

func (b *BentTrapezoid) SetCreationBounds(origin vectorpath.Point, size vectorpath.Point) {
	b.position = origin
	b.topWidth = size.P
	b.bottomWidth = size.P
	b.duration = size.T
}

func (b *BentTrapezoid) MarshalJSON() ([]byte, error) {
	var values = map[string]interface{}{
		"__TYPE__": "BentTrapezoid",
		"Shape": map[string]interface{}{
			"Position":     b.position,
			"TopWidth":     b.topWidth,
			"BottomWidth":  b.bottomWidth,
			"BottomOffset": b.bottomOffset,
			"Duration":     b.duration,
			"Bend":         b.bend,
		},
	}
	return json.Marshal(values)
}

func (b *BentTrapezoid) UnmarshalJSON(raw []byte) error {
	var values = make(map[string]interface{})
	err := json.Unmarshal(raw, &values)
	if err != nil {
		return err
	}

	position, ok := values["Position"]
	if !ok {
		return errors.New("bent trapezoid has missing key 'Position")
	}
	positionMap := position.(map[string]interface{})
	b.position.P = positionMap["P"].(float64)
	b.position.T = positionMap["T"].(float64)
	topWidth, ok := values["TopWidth"]
	if !ok {
		return errors.New("bent trapezoid has missing key 'TopWidth")
	}
	b.topWidth = topWidth.(float64)
	bottomWidth, ok := values["BottomWidth"]
	if !ok {
		return errors.New("bent trapezoid has missing key 'BottomWidth")
	}
	b.bottomWidth = bottomWidth.(float64)
	bottomOffset, ok := values["BottomOffset"]
	if !ok {
		return errors.New("bent trapezoid has missing key 'BottomOffset")
	}
	b.bottomOffset = bottomOffset.(float64)
	duration, ok := values["Duration"]
	if !ok {
		return errors.New("bent trapezoid has missing key 'Duration")
	}
	b.duration = duration.(float64)
	bend, ok := values["Bend"]
	if !ok {
		return errors.New("bent trapezoid has missing key 'Bend")
	}
	bendMap := bend.(map[string]interface{})
	b.bend.P = bendMap["P"].(float64)
	b.bend.T = bendMap["T"].(float64)

	return nil
}
