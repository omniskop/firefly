package main

import (
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
				},
				project.Element{
					ZIndex: 0,
					Shape:  shape.NewOrthogonalRectangle(vectorpath.Point{0.5, 5}, 0.5, 5),
				},
				project.Element{
					ZIndex: 0,
					Shape:  shape.NewOrthogonalRectangle(vectorpath.Point{0, 10}, 0.5, 5),
				},
				project.Element{
					ZIndex: 0,
					Shape:  shape.NewBentTrapezoid(vectorpath.Point{0.6, 15}, vectorpath.Point{0.4, 20}, 0.3, 0.5),
				},
			},
			Effects: []project.Effect{},
		},
		Audio: project.Audio{
			Title:  "song title",
			Author: "Salvatore Ganacci",
			Genres: []string{"trap", "house", "EDM"},
			File:   nil,
		},
		AudioOffset: 0,
	}

	editor.New(project)
	logrus.Info("Editor created")

	logrus.Info("Starting Application")

	app.Exec()
}
