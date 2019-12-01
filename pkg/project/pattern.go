package project

import (
	"image/color"

	"github.com/omniskop/firefly/pkg/project/vectorpath"
)

// Pattern describes how an element should be filled
type Pattern interface {
	Pattern() Pattern // this is just here to distinguish Pattern from an empty interface
}

// SolidColor fills an element with a solid color
type SolidColor struct {
	color.Color
}

// NewSolidColor returns a new SolidColor Pattern with the specified color
func NewSolidColor(c color.Color) *SolidColor {
	return &SolidColor{
		Color: c,
	}
}

// NewSolidColorRGBA returns a new SolidColor Pattern with the color specified as r, g, b and a components in the range of [0,255]
func NewSolidColorRGBA(r, g, b, a uint8) *SolidColor {
	return &SolidColor{
		Color: color.RGBA{R: r, G: g, B: b, A: a},
	}
}

// Pattern implements the Pattern interface
func (c *SolidColor) Pattern() Pattern {
	return c
}

// LinearGradient fills an element with a gradient between two points.
// The positions of the gradient are in local coordinates to the element,
// meaning that the top left position is (0,0) and the bottom right one is at (1,1).
type LinearGradient struct {
	Start GradientAnchorPoint
	Stop  GradientAnchorPoint
	Steps []GradientColorStep // the steps between the anchor points
}

// Pattern implements the Pattern interface
func (g *LinearGradient) Pattern() Pattern {
	return g
}

// A GradientAnchorPoint contains a position and the color that position should have in the gradient
type GradientAnchorPoint struct {
	color.Color
	vectorpath.Point
}

// GradientColorStep is a position on a gradient that has a specific color
type GradientColorStep struct {
	color.Color
	Position float64
}
