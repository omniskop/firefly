package scanner

import (
	"image/color"
	"math"
	"sort"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/omniskop/firefly/pkg/project/vectorpath"

	"github.com/omniskop/firefly/pkg/project"
)

// Frame contains the time of the frame and the pixel colors
type Frame struct {
	Time   float64
	Pixels []color.Color
}

// A Scanner can be used to scan lines of a project into separate pixel colors
type Scanner struct {
	scene   *project.Scene
	mapping *Mapping
	mutex   *sync.Mutex
}

// New creates a new scanner on the project. Size should be the number of led's.
func New(scene *project.Scene, size int) Scanner {
	return Scanner{
		scene:   scene,
		mapping: NewLinearMapping(size),
		mutex:   new(sync.Mutex),
	}
}

// GetPixelPosition returns the positions and the width of the pixel in the range of [0,1].
// It exposes the GetPixelPosition method of the mapping used by the scanner.
func (s Scanner) GetPixelPosition(p int) (float64, float64) {
	return s.mapping.GetPixelPosition(p)
}

func (s *Scanner) SetMapping(m Mapping) {
	s.mutex.Lock()
	s.mapping = &m
	s.mapping.fix()
	s.mutex.Unlock()
}

// Scans a line at the specified time and the returns the frame
func (s Scanner) Scan(time float64) Frame {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	size := s.mapping.Pixels()
	var frame = Frame{
		Time:   time,
		Pixels: make([]color.Color, size),
	}

	// set offsets to black
	for i := 0; i < s.mapping.StartOffset; i++ {
		frame.Pixels[i] = color.Black
	}
	for i := size - s.mapping.EndOffset; i < size; i++ {
		frame.Pixels[i] = color.Black
	}

	elements := s.scene.GetElementsAt(time)
	// logrus.WithField("elements", len(elements)).Debug("  ====== New Scan ======  ", time)

	// sort the elements in the correct ZIndex order
	sort.Slice(elements, func(i, j int) bool {
		return elements[i].ZIndex < elements[j].ZIndex
	})

	// Create a list of fragments. One for each element.
	// Theoretically an element could have more than one fragment but we currently ignore this case.
	// A fragment contains a start and a stop position.
	// Start is the first pixel where the elements starts to be visible
	// Stop is the last pixel where the elements is visible
	fragments := make([]struct {
		start float64
		stop  float64
	}, len(elements))

	for i, element := range elements {
		a, b := getPixelCoverageOfPath(element.Shape.Path(), time)
		fragments[i].start = a
		fragments[i].stop = b
	}

	// iterate through all pixels ...
	for pixelIndex := s.mapping.StartOffset; pixelIndex < size-s.mapping.EndOffset; pixelIndex++ {
		pixelPosition, pixelWidth := s.mapping.GetPixelPosition(pixelIndex)
		pixelPosition += pixelWidth / 2                             // use the center of the pixel for better results
		pixelInScene := vectorpath.Point{P: pixelPosition, T: time} // the location of the pixel in the scene
		var pixelColor color.Color = color.Black

		// ... and through all fragments ...
		for fragmentIndex, fragment := range fragments {
			// ... to check which fragments are visible in each pixel
			if fragment.start <= pixelPosition && fragment.stop >= pixelPosition {
				pixelColor = addColors(pixelColor, getFill(elements[fragmentIndex], pixelInScene))
			}
		}
		frame.Pixels[pixelIndex] = pixelColor
	}

	return frame
}

// getPixelCoverageOfPath returns the start and end positions [0, 1] where the shape is visible at the specific time
func getPixelCoverageOfPath(path vectorpath.Path, time float64) (float64, float64) {
	currentPoint := path.Start
	var edges []vectorpath.Point
	for _, segment := range path.Segments {
		newPoint := currentPoint.Add(segment.EndPoint())
		if (currentPoint.T < time && newPoint.T > time) || (currentPoint.T > time && newPoint.T < time) {
			edges = append(edges, vectorpath.Point{P: currentPoint.P + getSegmentEdge(segment, time-currentPoint.T), T: time})
		}
		currentPoint = newPoint
	}

	if len(edges) >= 2 {
		return math.Min(edges[0].P, edges[1].P), math.Max(edges[0].P, edges[1].P)
	}
	logrus.WithField("edges", len(edges)).Warn("GetPixelCoverage: not enough edges found")
	return 0, 0
}

// getSegmentEdge returns the position where the segment surpasses the point in time.
// time and the returned position are relative to the start of the segment
func getSegmentEdge(segment vectorpath.Segment, time float64) float64 {
	switch obj := segment.(type) {
	case *vectorpath.Line:
		return obj.P * (time / obj.T)
	case *vectorpath.QuadCurve:
		var foundProgress = math.NaN()
		if !floatsEqual(obj.Control.T*2, obj.End.T) {
			p1 := (obj.Control.T - math.Sqrt(math.Pow(obj.Control.T, 2)-2*obj.Control.T*time+obj.End.T*time)) / (2*obj.Control.T - obj.End.T)
			p2 := (obj.Control.T + math.Sqrt(math.Pow(obj.Control.T, 2)-2*obj.Control.T*time+obj.End.T*time)) / (2*obj.Control.T - obj.End.T)
			// theoretically both could be correct but we currently only support one hit
			if p1 >= 0 && p1 <= 1 {
				foundProgress = p1
			}
			if p2 >= 0 && p2 <= 1 {
				foundProgress = p2
			}
		} else {
			foundProgress = time / obj.End.T
		}

		if math.IsNaN(foundProgress) {
			return 0
		}
		return 2*foundProgress*(1-foundProgress)*obj.Control.P + math.Pow(foundProgress, 2)*obj.End.P
	case *vectorpath.CubicCurve:
		return 0
	}
	return 0
}

// getFill takes an element and a point inside it to return the correct color according to the pattern of the element.
func getFill(element *project.Element, point vectorpath.Point) color.Color {
	bounds := element.Shape.Bounds()
	point = point.Sub(bounds.Location)
	point = vectorpath.Point{
		P: point.P / bounds.Dimensions.P,
		T: point.T / bounds.Dimensions.T,
	}
	switch pattern := element.Pattern.(type) {
	case *project.SolidColor:
		return pattern.Color
	case *project.LinearGradient:
		toPoint := point.Sub(pattern.Start.Point)                // vector from the start of the gradient to the point of interest
		gradTrack := pattern.Stop.Point.Sub(pattern.Start.Point) // vector from start to end of the gradient

		// calculate the progress through projection
		progress := dotProduct(toPoint, gradTrack) / math.Pow(length(gradTrack), 2)

		return interpolateColors(pattern.Start.Color, pattern.Stop.Color, progress)
	default:
		return nil
	}
}

// interpolateColors interpolates linearly between colorA and colorB bases on progress.
// progress gets clamped between 0 and 1
func interpolateColors(colorA color.Color, colorB color.Color, progress float64) color.Color {
	progress = math.Min(1, math.Max(0, progress))
	aR, aG, aB, aA := colorToFloats(colorA)
	bR, bG, bB, bA := colorToFloats(colorB)
	return floatsToColor(
		aR+(bR-aR)*progress,
		aG+(bG-aG)*progress,
		aB+(bB-aB)*progress,
		aA+(bA-aA)*progress,
	)
}

// colorToFloats returns the r, g, b and a values of the color as floating point numbers between [0, 65536[
func colorToFloats(c color.Color) (float64, float64, float64, float64) {
	r, g, b, a := c.RGBA()
	return float64(r), float64(g), float64(b), float64(a)
}

// floatsToColors returns a new colors based on the rgba components as floats between [0, 65536[
func floatsToColor(r float64, g float64, b float64, a float64) color.Color {
	return color.RGBA64{
		R: uint16(math.Round(r)),
		G: uint16(math.Round(g)),
		B: uint16(math.Round(b)),
		A: uint16(math.Round(a)),
	}
}

// addColors combines the colors one and two in a way that two overlays one.
// The alpha value of two will be used to blend the color.
// If two has an alpha value of 255 it will completely override one and if the alpha value is 0 one will be fully visible.
// The transparency of one will be ignored and the resulting color is completely opaque.
func addColors(one color.Color, two color.Color) color.Color {
	oneR, oneG, oneB, _ := colorToFloats(one)
	twoR, twoG, twoB, twoA := colorToFloats(two)
	progress := twoA / 65535 // 65535 = 2^16-1 is the highest possible value for twoA
	blendR := oneR*(1-progress) + twoR*progress
	blendG := oneG*(1-progress) + twoG*progress
	blendB := oneB*(1-progress) + twoB*progress

	return floatsToColor(blendR, blendG, blendB, 65535)
}

func dotProduct(a, b vectorpath.Point) float64 {
	return a.P*b.P + a.T*b.T
}

func length(a vectorpath.Point) float64 {
	return math.Sqrt(a.P*a.P + a.T*a.T)
}

func floatsEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.00000001
}
