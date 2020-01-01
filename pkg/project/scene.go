package project

// Scene contains all the visual elements of a project
type Scene struct {
	Elements []Element
	Effects  []Effect
}

func (s Scene) GetElementsAt(time float64) []Element {
	var out []Element
	for _, element := range s.Elements {
		if element.Shape.Bounds().IncludesTime(time) {
			out = append(out, element)
		}
	}
	return out
}
