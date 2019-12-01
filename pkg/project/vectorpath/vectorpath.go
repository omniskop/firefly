// Package vectorpath contains definitions for shape paths
package vectorpath

import (
	"fmt"
	"math"
)

// A Point is a position in the scene
type Point struct {
	P float64 // P is the position of the point along the drawing axis
	T float64 // T is the point in time in seconds
}

// Sub returns the result of subtracting the parameter from this point
func (p Point) Sub(b Point) Point {
	return Point{
		T: p.T - b.T,
		P: p.P - b.P,
	}
}

// SubComponents returns the result of subtracting the parameters from this point
func (p Point) SubComponents(pos float64, time float64) Point {
	return Point{
		P: p.P - pos,
		T: p.T - time,
	}
}

// Add returns the result of adding this point with the parameter
func (p Point) Add(b Point) Point {
	return Point{
		T: p.T + b.T,
		P: p.P + b.P,
	}
}

// AddComponents returns the result of adding this point to the parameters
func (p Point) AddComponents(pos float64, time float64) Point {
	return Point{
		P: p.P + pos,
		T: p.T + time,
	}
}

// Invert returns the inverted form if this point
func (p Point) Invert() Point {
	return Point{
		T: -p.T,
		P: -p.P,
	}
}

// String implements the stringer interface
func (p Point) String() string {
	return fmt.Sprintf("{%.2f, %.2f}", p.P, p.T)
}

// Interpolate linearly between point a and b by the factor f
func Interpolate(a, b Point, f float64) Point {
	return Point{
		P: a.P + (b.P-a.P)*f,
		T: a.T + (b.T-a.T)*f,
	}
}

// Path contains segments that are positioned relative to P
type Path struct {
	P        float64
	Segments []Segment
}

// Duration returns the duration (length on the time axis) of the path
func (p Path) Duration() float64 {
	oldest := p.P
	for _, segment := range p.Segments {
		if segment.OldestPointInTime() > oldest {
			oldest = segment.OldestPointInTime()
		}
	}
	return oldest
}

// PointAfter returns the point where the path will be after drawing n segments starting at 'start'.
// If n is larger than the number of segments it will return the point after the full path has been drawn which is the same as the starting position.
func (p Path) PointAfter(start float64, n int) Point {
	point := Point{
		T: start,
		P: p.P,
	}
	if n >= len(p.Segments) {
		return point // shortcut because the path needs to be closed the end position of the path is the same as the starting one
	}
	for i := 0; i < n; i++ {
		point = point.Add(p.Segments[i].EndPoint())
	}
	return point
}

// A Segment is a part of a larger path
type Segment interface {
	Move(Point)
	EndPoint() Point
	OldestPointInTime() float64 // returns the oldest point in time of the segment
}

// A Line to a point
type Line struct {
	Point
}

// Move moves the Line by some amount
func (l *Line) Move(diff Point) {
	l.Point = l.Point.Add(diff)
}

// EndPoint of the line
func (l *Line) EndPoint() Point {
	return l.Point
}

// OldestPointInTime returns the end point of the line
func (l *Line) OldestPointInTime() float64 {
	return l.T
}

// A QuadCurve from the last point with control and end points
type QuadCurve struct {
	Control Point
	End     Point
}

// Move moves the QuadCurve by some amount
func (curve *QuadCurve) Move(diff Point) {
	curve.Control = curve.Control.Add(diff)
	curve.End = curve.End.Add(diff)
}

// EndPoint of the QuadCurve
func (curve *QuadCurve) EndPoint() Point {
	return curve.End
}

// OldestPointInTime returns the oldest point of the curve
func (curve *QuadCurve) OldestPointInTime() float64 {
	if curve.Control.T > curve.End.T {
		return curve.Control.T
	}
	return curve.End.T
}

// A CubicCurve from the last point with two control points and an end point
type CubicCurve struct {
	ControlA Point
	ControlB Point
	End      Point
}

// Move moves the CubicCurve by some amount
func (curve *CubicCurve) Move(diff Point) {
	curve.ControlA = curve.ControlA.Add(diff)
	curve.ControlB = curve.ControlB.Add(diff)
	curve.End = curve.End.Add(diff)
}

// EndPoint of the CubicCurve
func (curve *CubicCurve) EndPoint() Point {
	return curve.End
}

// OldestPointInTime returns the oldest point of the curve
func (curve *CubicCurve) OldestPointInTime() float64 {
	if curve.ControlA.T >= curve.ControlB.T {
		if curve.ControlA.T >= curve.End.T {
			return curve.ControlA.T
		}
	} else {
		if curve.ControlB.T >= curve.End.T {
			return curve.ControlB.T
		}
	}
	return curve.End.T
}

// Clamp the value between min and max
func Clamp(value float64, min float64, max float64) float64 {
	// using min and max is faster than using an if
	return math.Min(max, math.Max(min, value))
}
