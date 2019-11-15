// Package project contains the firefly project structure
package project

// A Moment represents a point in time in a project in seconds
type Moment float64

// Project contains everything needed to describe a firefly project
type Project struct {
	Title          string            // title of the project
	Author         string            // name of the project author
	Tags           []string          // tags of the project
	AdditionalInfo map[string]string // additional information about the project for future extensibility
	Duration       float64           // the duration of the project in seconds
	Scene          Scene             // the visual elements of the project
	Audio          Audio             // the audio of the project
	AudioOffset    float64           // the offset of the audio timeline from the visual timeline. This can be negative
}
