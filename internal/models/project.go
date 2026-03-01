// Package models provides project-level data structures.
package models

import (
	"os"
	"path/filepath"
)

// Project holds metadata about the Go project being inspected.
type Project struct {
	Dir  string // Absolute path to the project root
	Name string // Derived short name (directory basename)
}

// DetectProject builds a Project from the current working directory.
func DetectProject() (Project, error) {
	dir, err := os.Getwd()
	if err != nil {
		return Project{}, err
	}
	name := filepath.Base(dir)
	return Project{Dir: dir, Name: name}, nil
}
