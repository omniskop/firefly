package project

// Element node is a single node inside of an ElementCollection
type ElementNode struct {
	self     *Element
	previous *ElementNode
	next     *ElementNode
}

func (en *ElementNode) Next() *ElementNode {
	return en.next
}

func (en *ElementNode) Previous() *ElementNode {
	return en.previous
}

func (en *ElementNode) Element() *Element {
	return en.self
}

// ElementCollection is a linked list that contains an ordered set of elements.
// Modifications of this list are *not* thread safe.
type ElementCollection struct {
	first *ElementNode
	last  *ElementNode
	size  int
}

func NewEmptyElementCollection() *ElementCollection {
	return &ElementCollection{nil, nil, 0}
}

func NewElementCollection(elements []*Element) *ElementCollection {
	col := &ElementCollection{size: len(elements)}
	if len(elements) == 0 {
		return col
	}
	col.first = &ElementNode{self: elements[0]}
	prev := col.first
	for _, element := range elements {
		prev.next = &ElementNode{self: element, previous: prev}
		prev = prev.next
	}
	col.last = prev
	return col
}

// Size of the collection
func (ec *ElementCollection) Size() int {
	return ec.size
}

// IsEmpty returns true if the collection is empty
func (ec *ElementCollection) IsEmpty() bool {
	return ec.size == 0
}

// DeleteAll elements from the collection
func (ec *ElementCollection) DeleteAll() {
	ec.first = nil
	ec.last = nil
	ec.size = 0
}

// At returns the ElementNode that the given index. As the collection is a linked list this function is potentially slow.
// If the index is inside the list it will return the node and true,
// if it is not nil and false will be returned.
func (ec *ElementCollection) nodeAt(i int) (*ElementNode, bool) {
	if i >= ec.size {
		// the index is out of bounds
		return nil, false
	}

	// TODO: traverse list from the back if the index is closer to the end
	prev := ec.first
	for c := 0; c < i; c++ {
		prev = prev.next
		if prev == nil {
			// the index is out of bounds
			return nil, false
		}
	}
	return prev, true
}

// At returns the element that the given index. As the collection is a linked list this function is potentially slow.
// If the index is inside the list it will return the stored element (which could be nil) and true,
// if it is not nil and false will be returned.
func (ec *ElementCollection) At(i int) (*Element, bool) {
	node, ok := ec.nodeAt(i)
	if !ok {
		return nil, false
	}
	return node.self, ok
}

// DeleteIndex will delete the element at the given index from the list.
// As the collection is a linked list this function is potentially slow.
// It returns false when the index is outside of the list and true if otherwise.
func (ec *ElementCollection) DeleteIndex(i int) bool {
	node, ok := ec.nodeAt(i)
	if !ok {
		return false
	}

	if node.previous != nil {
		node.previous.next = node.next
	}
	if node.next != nil {
		node.previous.previous = node.previous
	}

	if ec.first == node {
		ec.first = node.next
	}
	if ec.last == node {
		ec.last = node.previous
	}

	ec.size--
	return true
}

// IndexOf searches the given element in the collection.
// The boolean will indicate weather or not the element has been found.
// As the collection is a linked list this function is potentially slow.
func (ec *ElementCollection) IndexOf(element *Element) (int, bool) {
	node := ec.first
	for i := 0; node != nil; i++ {
		if node.self == element {
			return i, true
		}
		node = node.next
	}
	return -1, false
}

// InsertAt will insert the element in the collection at the given index.
// If the index would be outside the collection the element will not be inserted the function returns false.
// As the collection is a linked list this function is potentially slow.
func (ec *ElementCollection) InsertAt(element *Element, i int) bool {
	if ec.size == i {
		// the element will be appended to the end

		if ec.size == 0 {
			// this will be the first element
			node := &ElementNode{self: element, previous: nil, next: nil}
			ec.first = node
			ec.last = node
			ec.size++
			return true
		}

		node := &ElementNode{self: element, previous: ec.last, next: nil}
		ec.last.next = node
		ec.last = node
		ec.size++
		return true
	}

	newNextNode, ok := ec.nodeAt(i)
	if !ok {
		return false
	}
	prev := newNextNode.previous
	node := &ElementNode{self: element, previous: prev, next: newNextNode}
	prev.next = node
	newNextNode.previous = node

	ec.size++
	return true
}

// Append the element to the collection
func (ec *ElementCollection) Append(element *Element) {
	node := &ElementNode{self: element, previous: ec.last, next: nil}
	if ec.size == 0 {
		ec.first = node
		ec.last = node
		return
	}
	ec.last.next = node
	ec.last = node
}

// AppendElements appends multiple elements to the collection.
// The order will be preserved.
func (ec *ElementCollection) AppendElements(elements []*Element) {
	for _, element := range elements {
		ec.Append(element)
	}
}

// Prepend the element to the collection
func (ec *ElementCollection) Prepend(element *Element) {
	node := &ElementNode{self: element, previous: nil, next: ec.first}
	if ec.size == 0 {
		ec.first = node
		ec.last = node
		return
	}
	ec.first.previous = node
	ec.first = node
}

// PrependElements prepends multiple elements to the collection.
// The order will be preserved.
func (ec *ElementCollection) PrependElements(elements []*Element) {
	for _, element := range elements {
		ec.Prepend(element)
	}
}

// Elements returns a slice with all element inside the collection
func (ec *ElementCollection) Elements() []*Element {
	var out = make([]*Element, ec.size)
	node := ec.first
	for i := 0; node != nil; i++ {
		out[i] = node.self
		node = node.next
	}
	return out
}

// FirstNode returns the first node of the list that can be used to iterate over all Elements.
func (ec *ElementCollection) FirstNode() *ElementNode {
	return ec.first
}

// LastNode returns the last node of the list that can be used to iterate over all Elements.
func (ec *ElementCollection) LastNode() *ElementNode {
	return ec.last
}
