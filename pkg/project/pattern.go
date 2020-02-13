package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"image/color"

	"github.com/omniskop/firefly/pkg/project/vectorpath"
)

func MarshalColor(c color.Color) []byte {
	r, g, b, a := c.RGBA()
	return []byte(fmt.Sprintf(`{"R":%d,"G":%d,"B":%d,"A":%d}`, r, g, b, a))
}

func UnmarshalColor(raw []byte) (color.Color, error) {
	var c color.RGBA64
	err := json.Unmarshal(raw, &c)
	return c, err
}

// Pattern describes how an element should be filled
type Pattern interface {
	json.Marshaler

	Pattern() Pattern // this is just here to distinguish Pattern from an empty interface
	MirrorP()         // mirrors the pattern on the P axis
	Copy() Pattern
}

func UnmarshalPattern(raw []byte) (Pattern, error) {
	values := make(map[string]*json.RawMessage)
	err := json.Unmarshal(raw, &values)
	if err != nil {
		return nil, err
	}

	rawType, ok := values["__TYPE__"]
	if !ok {
		return nil, fmt.Errorf("pattern has missing key 'Type'")
	}
	var patternType string
	err = json.Unmarshal(*rawType, &patternType)
	if err != nil {
		return nil, fmt.Errorf("pattern has invalid key 'Type'")
	}
	if _, ok := values["Pattern"]; !ok {
		return nil, fmt.Errorf("pattern has missing key 'Pattern'")
	}

	switch patternType {
	case "SolidColor":
		data := new(SolidColor)
		err = json.Unmarshal(*values["Pattern"], data)
		return data, err
	case "LinearGradient":
		data := new(LinearGradient)
		err = json.Unmarshal(*values["Pattern"], data)
		return data, err
	default:
		// I considered using a json.UnmarshalTypeError but decided against it because it has a bunch of field that
		// i would not fill and it would probably end up less descriptive than just a simple error.
		return nil, fmt.Errorf("pattern has unknown type %q", *values["__TYPE__"])
	}
}

// SolidColor fills an element with a solid color
type SolidColor struct {
	color.Color
}

var _ Pattern = (*SolidColor)(nil) // make sure SolidColor implements the Pattern interface

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

func (c *SolidColor) MirrorP() {
	// no action needed
}

func (c *SolidColor) Copy() Pattern {
	r, g, b, a := c.RGBA()
	return NewSolidColor(color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)})
}

func (c *SolidColor) MarshalJSON() ([]byte, error) {
	var values = map[string]interface{}{
		"__TYPE__": "SolidColor",
		"Pattern": map[string]json.RawMessage{
			"Color": MarshalColor(c.Color),
		},
	}
	return json.Marshal(values)
}

func (c *SolidColor) UnmarshalJSON(raw []byte) error {
	var values = make(map[string]*json.RawMessage)
	err := json.Unmarshal(raw, &values)
	if err != nil {
		return err
	}

	rawColor, ok := values["Color"]
	if !ok {
		return errors.New("solid color has missing key 'Color")
	}

	c.Color, err = UnmarshalColor(*rawColor)
	return err
}

// LinearGradient fills an element with a gradient between two points.
// The positions of the gradient are in local coordinates to the element,
// meaning that the top left position is (0,0) and the bottom right one is at (1,1).
type LinearGradient struct {
	Start GradientAnchorPoint
	Stop  GradientAnchorPoint
	Steps []GradientColorStep // the steps between the anchor points
}

var _ Pattern = (*LinearGradient)(nil) // make sure LinearGradient implements the Pattern interface

// NewLinearGradient creates a new LinearGradient with the given start and stop colors
func NewLinearGradient(a color.Color, b color.Color) *LinearGradient {
	return &LinearGradient{
		Start: GradientAnchorPoint{
			Color: a,
			Point: vectorpath.Point{P: 0.5, T: 0},
		},
		Stop: GradientAnchorPoint{
			Color: b,
			Point: vectorpath.Point{P: 0.5, T: 1},
		},
		Steps: nil,
	}
}

// Pattern implements the Pattern interface
func (g *LinearGradient) Pattern() Pattern {
	return g
}

func (g *LinearGradient) MirrorP() {
	g.Start.Point.P = 1 - g.Start.Point.P
	g.Stop.Point.P = 1 - g.Stop.Point.P
}

func (g *LinearGradient) Copy() Pattern {
	return &LinearGradient{
		Start: g.Start.Copy(),
		Stop:  g.Stop.Copy(),
		Steps: nil,
	}
}

func (g *LinearGradient) MarshalJSON() ([]byte, error) {
	var values = map[string]interface{}{
		"__TYPE__": "LinearGradient",
		"Pattern": map[string]interface{}{
			"Start": g.Start,
			"Stop":  g.Stop,
			"Steps": g.Steps,
		},
	}
	return json.Marshal(values)
}

// A GradientAnchorPoint contains a position and the color that position should have in the gradient
type GradientAnchorPoint struct {
	color.Color
	vectorpath.Point
}

func (p GradientAnchorPoint) Copy() GradientAnchorPoint {
	r, g, b, a := p.Color.RGBA()
	return GradientAnchorPoint{
		Color: color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)},
		Point: p.Point,
	}
}

func (p GradientAnchorPoint) MarshalJSON() ([]byte, error) {
	var values = map[string]interface{}{
		"Color": json.RawMessage(MarshalColor(p.Color)),
		"Point": p.Point,
	}
	return json.Marshal(values)
}

func (p *GradientAnchorPoint) UnmarshalJSON(raw []byte) error {
	var values = make(map[string]*json.RawMessage)
	err := json.Unmarshal(raw, &values)
	if err != nil {
		return err
	}

	point, ok := values["Point"]
	if !ok {
		return errors.New("gradient anchor point has missing key 'Point'")
	}
	color, ok := values["Color"]
	if !ok {
		return errors.New("gradient anchor point has missing key 'Color'")
	}

	err = json.Unmarshal(*point, &p.Point)
	if err != nil {
		return err
	}
	p.Color, err = UnmarshalColor(*color)
	if err != nil {
		return err
	}

	return nil
}

// GradientColorStep is a position on a gradient that has a specific color
type GradientColorStep struct {
	color.Color
	Position float64
}

func (s GradientColorStep) MarshalJSON() ([]byte, error) {
	var values = map[string]interface{}{
		"Color":    json.RawMessage(MarshalColor(s.Color)),
		"Position": s.Position,
	}
	return json.Marshal(values)
}

func (s *GradientColorStep) UnmarshalJSON(raw []byte) error {
	var values = make(map[string]*json.RawMessage)
	err := json.Unmarshal(raw, &values)
	if err != nil {
		return err
	}

	position, ok := values["Position"]
	if !ok {
		return errors.New("gradient anchor point has missing key 'Position'")
	}
	color, ok := values["Color"]
	if !ok {
		return errors.New("gradient anchor point has missing key 'Color'")
	}

	err = json.Unmarshal(*position, &s.Position)
	if err != nil {
		return err
	}
	s.Color, err = UnmarshalColor(*color)
	if err != nil {
		return err
	}

	return nil
}
