// Go Developer Dashboard - main entrypoint.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/chmenegatti/devdash/internal/app"
	"github.com/chmenegatti/devdash/internal/logs"
	"github.com/chmenegatti/devdash/internal/models"
	"github.com/chmenegatti/devdash/internal/state"
)

func main() {
	project, err := models.DetectProject()
	if err != nil {
		logs.Errorf("project detection failed: %v", err)
		fmt.Fprintf(os.Stderr, "Error detecting project: %v\n", err)
		os.Exit(1)
	}

	if err := logs.SetFile(filepath.Join(project.Dir, ".devdash.log")); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not configure log file: %v\n", err)
	} else {
		logs.Infof("logger initialized at %s", logs.FilePath())
	}

	ds := state.New(project.Dir, project.Name, resolveVersion())
	m := app.New(ds)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logs.Errorf("dashboard runtime failed: %v", err)
		fmt.Fprintf(os.Stderr, "Error running dashboard: %v\n", err)
		os.Exit(1)
	}
}

func resolveVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}

	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}

	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" && setting.Value != "" {
			if len(setting.Value) > 7 {
				return "dev-" + setting.Value[:7]
			}
			return "dev-" + setting.Value
		}
	}

	return "dev"
}
