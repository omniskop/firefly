package main

import (
	"image/color"
	"runtime"

	"github.com/omniskop/firefly/pkg/project/shape"

	"github.com/omniskop/firefly/cmd/firefly/editor"

	"github.com/omniskop/firefly/pkg/project"
	"github.com/omniskop/firefly/pkg/project/vectorpath"
	"github.com/sirupsen/logrus"
	"github.com/therecipe/qt/widgets"
)

var app *widgets.QApplication

func init() {
	// not shure if this is realy needed
	runtime.LockOSThread()
}

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.Info("FireFly starting...")

	app = createApplication()

	logrus.Info("Application created")

	project := &project.Project{
		Title:          "test title",
		Author:         "omniskop",
		Tags:           []string{"hot", "awesome"},
		AdditionalInfo: map[string]string{"demo": "true"},
		Duration:       300, // 5 Minutes
		Scene: project.Scene{
			Elements: []project.Element{
				project.Element{
					ZIndex: 0,
					Shape:  shape.NewOrthogonalRectangle(vectorpath.Point{0, 0}, 0.5, 5),
					Pattern: &project.LinearGradient{
						Start: project.GradientAnchorPoint{
							Color: color.RGBA{R: 255, G: 0, B: 0, A: 255},
							Point: vectorpath.Point{P: 0, T: 0},
						},
						Stop: project.GradientAnchorPoint{
							Color: color.RGBA{R: 0, G: 0, B: 255, A: 255},
							Point: vectorpath.Point{P: 1, T: 0},
						},
						Steps: nil,
					},
				},
				project.Element{
					ZIndex:  0,
					Shape:   shape.NewOrthogonalRectangle(vectorpath.Point{0.5, 5}, 0.5, 5),
					Pattern: project.NewSolidColorRGBA(255, 0, 0, 255),
				},
				project.Element{
					ZIndex:  0,
					Shape:   shape.NewOrthogonalRectangle(vectorpath.Point{0, 10}, 0.5, 5),
					Pattern: project.NewSolidColorRGBA(255, 0, 0, 255),
				},
				project.Element{
					ZIndex: 0,
					Shape:  shape.NewBentTrapezoid(vectorpath.Point{0.6, 15}, vectorpath.Point{0.4, 20}, 0.3, 0.5),
					Pattern: &project.LinearGradient{
						Start: project.GradientAnchorPoint{
							Color: color.RGBA{R: 255, G: 0, B: 129, A: 255},
							Point: vectorpath.Point{P: 0, T: 0},
						},
						Stop: project.GradientAnchorPoint{
							Color: color.RGBA{R: 255, G: 99, B: 0, A: 255},
							Point: vectorpath.Point{P: 1, T: 1},
						},
						Steps: nil,
					},
				},
				project.Element{
					ZIndex: 0,
					Shape:  shape.NewBentTrapezoid(vectorpath.Point{0.5, 20}, vectorpath.Point{0.5, 25}, 0.5, 0.5),
					Pattern: &project.LinearGradient{
						Start: project.GradientAnchorPoint{
							Color: color.RGBA{R: 255, G: 0, B: 0, A: 255},
							Point: vectorpath.Point{P: 0, T: 0},
						},
						Stop: project.GradientAnchorPoint{
							Color: color.RGBA{R: 0, G: 0, B: 255, A: 255},
							Point: vectorpath.Point{P: 1, T: 0},
						},
						Steps: nil,
					},
				},
			},
			Effects: []project.Effect{},
		},
		Audio: project.Audio{
			//Title:  "Zeitansage",
			Title:  "All That Matters",
			Author: "Professor Kliq",
			Genres: []string{"electronic", "dance"},
			File:   nil,
		},
		AudioOffset: 0,
	}

	editor.New(project, ApplicationCallbacks)
	logrus.Info("Editor created")

	logrus.Info("Starting Application")

	app.Exec()
}
