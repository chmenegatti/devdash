package modules

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chmenegatti/devdash/internal/services"
	"github.com/chmenegatti/devdash/internal/state"
)

// RunBinarySize builds the project into a temp binary, measures its size, and cleans up.
// Blocking — call from a tea.Cmd.
func RunBinarySize(projectDir string) state.BinaryResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Resolve to absolute path so the temp binary path is reliable
	absDir, err := filepath.Abs(projectDir)
	if err != nil {
		return state.BinaryResult{
			Status: state.StatusError,
			Err:    "resolve path: " + err.Error(),
		}
	}

	tmpBin := filepath.Join(absDir, ".devdash_tmp_binary")
	defer func() {
		if err := os.Remove(tmpBin); err != nil && !os.IsNotExist(err) {
			_ = err
		}
	}()

	res := services.RunCommand(ctx, absDir, "go", "build", "-o", tmpBin, ".")
	if res.Err != nil {
		target, err := detectMainBuildTarget(ctx, absDir)
		if err != nil {
			return state.BinaryResult{
				Status: state.StatusError,
				Err:    res.Err.Error() + "\n" + res.Stderr + "\n" + err.Error(),
			}
		}
		res = services.RunCommand(ctx, absDir, "go", "build", "-o", tmpBin, target)
	}

	if res.Err != nil {
		return state.BinaryResult{
			Status: state.StatusError,
			Err:    res.Err.Error() + "\n" + res.Stderr,
		}
	}

	info, err := os.Stat(tmpBin)
	if err != nil {
		return state.BinaryResult{
			Status: state.StatusError,
			Err:    "stat failed: " + err.Error(),
		}
	}

	return state.BinaryResult{
		Status: state.StatusDone,
		Size:   info.Size(),
	}
}

func detectMainBuildTarget(ctx context.Context, projectDir string) (string, error) {
	list := services.RunCommand(ctx, projectDir, "go", "list", "-f", "{{if eq .Name \"main\"}}{{.ImportPath}}{{end}}", "./...")
	if list.Err != nil {
		return "", fmt.Errorf("detect main package: %w\n%s", list.Err, strings.TrimSpace(list.Stderr))
	}

	for _, line := range strings.Split(list.Stdout, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			return line, nil
		}
	}

	return "", fmt.Errorf("no main package found in project")
}
