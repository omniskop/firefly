package editor

import (
	"image/color"

	"github.com/therecipe/qt/core"

	"github.com/omniskop/firefly/pkg/project"
	"github.com/therecipe/qt/gui"
)

func NewQColorFromColor(col color.Color) *gui.QColor {
	r, g, b, a := col.RGBA()
	// the color values are between 0 and 0xffff
	// by dividing them by 257 they are converted to a range between 0 and 255
	// (2^16 - 1) / 257 = 255
	return gui.NewQColor3(int(r/257), int(g/257), int(b/257), int(a/257))
}

func NewColorFromQColor(qColor *gui.QColor) color.Color {
	// TODO: check if qColor.GetRgb() would work
	return color.RGBA{
		R: uint8(qColor.Red()),
		G: uint8(qColor.Green()),
		B: uint8(qColor.Blue()),
		A: uint8(qColor.Alpha()),
	}
}

func NewQLinearGradientFromLinearGradient(grad *project.LinearGradient) *gui.QLinearGradient {
	// we can't use qtPoint for the conversion here because these are not scene coordinates
	var qgradient *gui.QLinearGradient
	if verticalTimeAxis {
		qgradient = gui.NewQLinearGradient3(grad.Start.Point.P, grad.Start.Point.T, grad.Stop.Point.P, grad.Stop.Point.T)
	} else {
		qgradient = gui.NewQLinearGradient3(grad.Start.Point.T, grad.Start.Point.P, grad.Stop.Point.T, grad.Stop.Point.P)
	}
	qgradient.SetCoordinateMode(gui.QGradient__ObjectMode) // object mode => (0,0) <-> (1,1)
	qgradient.SetColorAt(0, NewQColorFromColor(grad.Start))
	for _, step := range grad.Steps {
		qgradient.SetColorAt(step.Position, NewQColorFromColor(step.Color))
	}
	qgradient.SetColorAt(1, NewQColorFromColor(grad.Stop))
	return qgradient
}

func NewQBrushFromPattern(pat project.Pattern) *gui.QBrush {
	switch cast := pat.(type) {
	case *project.SolidColor:
		return gui.NewQBrush3(NewQColorFromColor(cast), core.Qt__SolidPattern)
	case *project.LinearGradient:
		return gui.NewQBrush10(NewQLinearGradientFromLinearGradient(cast))
	default:
		return gui.NewQBrush3(gui.NewQColor3(240, 107, 255, 255), core.Qt__SolidPattern)
	}
}
