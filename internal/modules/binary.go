package modules

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/cesar/devdash/internal/services"
	"github.com/cesar/devdash/internal/state"
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

	res := services.RunCommand(ctx, absDir, "go", "build", "-o", tmpBin, "./cmd/dashboard")

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
