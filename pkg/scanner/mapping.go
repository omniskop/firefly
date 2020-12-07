package scanner

type Mapping struct {
	Reversed    bool
	StartOffset int
	EndOffset   int
	Segments    []Segment
}

func NewLinearMapping(pixels int) *Mapping {
	return &Mapping{
		Reversed:    false,
		StartOffset: 0,
		EndOffset:   0,
		Segments:    []Segment{{PixelSize: pixels, To: 1}},
	}
}

func (m *Mapping) Pixels() int {
	sum := m.StartOffset + m.EndOffset
	for _, s := range m.Segments {
		sum += s.PixelSize
	}
	return sum
}

func (m *Mapping) AddStop(position float64, pixels int) {
	newSegments := make([]Segment, 0, len(m.Segments))
	didAdd := false
	for _, seg := range m.Segments {
		if seg.To > position {
			newSegments = append(newSegments, Segment{pixels, position})
			didAdd = true
		}
		newSegments = append(newSegments, seg)
	}
	if !didAdd {
		newSegments = append(newSegments, Segment{pixels, position})
	}
	m.Segments = newSegments
}

// GetPixelPosition returns the positions and the width of the pixel in the range of [0,1]
func (m *Mapping) GetPixelPosition(p int) (float64, float64) {
	pixelOffset := m.StartOffset
	var lastPosition float64
	// iterate through the segments
	for _, seg := range m.Segments {
		// the first pixel of the segment is currently in pixelOffset
		// when we add the PixelSize of the segment we get the last pixel (or first pixel of next segment)
		if pixelOffset+seg.PixelSize > p {
			// we found the segment that this pixel is in
			p := float64(p-pixelOffset) / float64(seg.PixelSize)                // [0,1] position of pixel in segment
			pixelWidth := float64(seg.To-lastPosition) / float64(seg.PixelSize) // "width" of a pixel in position space
			// calculate position of the pixel
			return lastPosition + (seg.To-lastPosition)*p, pixelWidth
		}
		pixelOffset += seg.PixelSize
		lastPosition = seg.To
	}
	return 0, 0
}

// fix makes sure that the mapping is valid
func (m *Mapping) fix() {
	if len(m.Segments) == 0 {
		m.Segments = []Segment{{PixelSize: 1, To: 1}}
	} else {
		// make sure that the last segment ends at 1.0
		m.Segments[len(m.Segments)-1].To = 1
	}
}

type Segment struct {
	PixelSize int
	To        float64
}
