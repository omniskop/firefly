package editor

import (
	"github.com/omniskop/firefly/pkg/project"
)

// elementList contains a modifiable list of graphicItems
// onChange will be called whenever the list changes and could potentially be made optional if that is required
type elementList struct {
	elements []*elementGraphicsItem
	onChange func()
}

func (s elementList) isEmpty() bool {
	return len(s.elements) == 0
}

func (s elementList) contains(searched *elementGraphicsItem) bool {
	for _, element := range s.elements {
		if element == searched {
			return true
		}
	}
	return false
}

func (s elementList) copyElements() []*project.Element {
	out := make([]*project.Element, len(s.elements))

	for i, item := range s.elements {
		out[i] = item.element.Copy()
	}
	return out
}

func (s *elementList) removeIfFound(searched *elementGraphicsItem) {
	for i, element := range s.elements {
		if element == searched {
			element.hideHandles()
			s.elements[i] = s.elements[len(s.elements)-1]
			s.elements = s.elements[:len(s.elements)-2]
		}
	}
	s.onChange()
}

func (s *elementList) clear() {
	for _, element := range s.elements {
		element.hideHandles()
	}
	s.elements = []*elementGraphicsItem{}
	s.onChange()
}

func (s *elementList) add(item *elementGraphicsItem) {
	item.showHandles()
	s.elements = append(s.elements, item)
	s.onChange()
}

func (s *elementList) set(item *elementGraphicsItem) {
	s.clear()
	s.add(item)
	s.onChange()
}
