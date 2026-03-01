// Go Developer Dashboard - main entrypoint.
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/chmenegatti/devdash/internal/app"
	"github.com/chmenegatti/devdash/internal/models"
	"github.com/chmenegatti/devdash/internal/state"
)

func main() {
	project, err := models.DetectProject()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting project: %v\n", err)
		os.Exit(1)
	}

	ds := state.New(project.Dir, project.Name)
	m := app.New(ds)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running dashboard: %v\n", err)
		os.Exit(1)
	}
}
