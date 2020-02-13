package shape

import (
	"encoding/json"
	"fmt"

	"github.com/omniskop/firefly/pkg/project/vectorpath"
)

// Shape is the visual part of an element on the scene
type Shape interface {
	json.Marshaler

	Time() float64              // point in time when the shape start
	Duration() float64          // duration of time that the shape takes up
	Width() float64             // visual width of the shape
	Bounds() vectorpath.Rect    // outer rectangular bounds of the shape
	Move(vectorpath.Point)      // move the shape by some amount
	Origin() vectorpath.Point   // get the point where the path of the shape starts (does not have to be the same as Bounds().Location)
	SetOrigin(vectorpath.Point) // set the origin of the shape

	Path() vectorpath.Path
	Handles() []vectorpath.Point                          // returns all points where the user can manipulate the shape
	SetHandle(int, vectorpath.Point)                      // set new position of a handle
	SetCreationBounds(vectorpath.Point, vectorpath.Point) // sets the size of the shape in an intuitive way for the user
	MirrorP()                                             // mirrors the shape on the P axis

	Copy() Shape // creates a deep copy of the shape
}

func Unmarshal(raw []byte) (Shape, error) {
	values := make(map[string]*json.RawMessage)
	err := json.Unmarshal(raw, &values)
	if err != nil {
		return nil, err
	}

	rawType, ok := values["__TYPE__"]
	if !ok {
		return nil, fmt.Errorf("shape has missing key 'Type'")
	}
	var shapeType string
	err = json.Unmarshal(*rawType, &shapeType)
	if err != nil {
		return nil, fmt.Errorf("shape has invalid key 'Type'")
	}
	if _, ok := values["Shape"]; !ok {
		return nil, fmt.Errorf("shape has missing key 'Shape'")
	}

	switch shapeType {
	case "OrthogonalRectangle":
		data := new(OrthogonalRectangle)
		err = json.Unmarshal(*values["Shape"], data)
		return data, err
	case "BentTrapezoid":
		data := new(BentTrapezoid)
		err = json.Unmarshal(*values["Shape"], data)
		return data, err
	default:
		// I considered using a json.UnmarshalTypeError but decided against it because it has a bunch of field that
		// i would not fill and it would probably end up less descriptive than just a simple error.
		return nil, fmt.Errorf("shape has unknown type %q", *values["__TYPE__"])
	}
}
