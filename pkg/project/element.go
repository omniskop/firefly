package project

// An Element that is located at a specific point in time in the scene
type Element struct {
	ZIndex float64 // a coordinate relative to other elements in the scene. Higher numbers will be drawn ontop of lower ones
	Shape  Shape   // the actual visual shape of the element
}
