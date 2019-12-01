package shape

import (
	"math"

	"github.com/omniskop/firefly/pkg/project"
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

var _ project.Shape = (*BentTrapezoid)(nil) // make sure BentTrapezoid implements the Shape interface

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

func (b *BentTrapezoid) Time() float64 {
	return b.position.T
}

func (b *BentTrapezoid) Duration() float64 {
	return b.duration
}

func (b *BentTrapezoid) Width() float64 {
	return math.Max(b.topWidth, b.bottomOffset+b.bottomWidth)
}

func (b *BentTrapezoid) Path() vectorpath.Path {
	return vectorpath.Path{
		P: b.position.P,
		Segments: []vectorpath.Segment{
			&vectorpath.Line{Point: vectorpath.Point{P: b.topWidth, T: 0}}, // top right
			//&vectorpath.Line{Point: vectorpath.Point{ // bottom right
			//	P: (b.bottomOffset + b.bottomWidth) - (b.topWidth),
			//	T: b.duration,
			//}},
			&vectorpath.QuadCurve{ // bottom right
				Control: vectorpath.Point{
					P: b.test((b.bottomOffset+b.bottomWidth)-b.topWidth, b.bend.P),
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
					P: -b.test(b.bottomOffset, 1-b.bend.P),
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

//func (b *BentTrapezoid) interpolateFromLeftToRight(a, b, p float64) {
//	if a < b {
//		return
//	}
//}

func (b *BentTrapezoid) test(a, p float64) float64 {
	if a > 0 {
		return a * p
	} else {
		return a * (1 - p)
	}
}

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

func (b *BentTrapezoid) Move(offset vectorpath.Point) {
	b.position = b.position.Add(offset)
}

func (b *BentTrapezoid) Location() vectorpath.Point {
	return b.position
}

func (b *BentTrapezoid) SetLocation(point vectorpath.Point) {
	b.position = point
}
