// Package storage contains everything needed to save and load a project
package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/omniskop/firefly/pkg/project"
)

func LoadFile(filename string) (*project.Project, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return Load(file)
}

func Load(input io.Reader) (*project.Project, error) {
	decoder := json.NewDecoder(input)
	var proj = new(project.Project)
	err := decoder.Decode(proj)
	if err != nil {
		return proj, fmt.Errorf("couldn't load project file: %w", err)
	}
	return proj, nil
}

func SaveFile(filename string, proj *project.Project) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return Save(file, proj)
}

func Save(output io.Writer, proj *project.Project) error {
	encoder := json.NewEncoder(output)
	return encoder.Encode(proj)
}
