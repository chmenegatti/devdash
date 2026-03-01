package logs

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	mu      sync.Mutex
	logPath = filepath.Join(os.TempDir(), "devdash.log")
)

func SetFile(path string) error {
	if path == "" {
		return fmt.Errorf("log path cannot be empty")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()
	logPath = path
	return nil
}

func FilePath() string {
	mu.Lock()
	defer mu.Unlock()
	return logPath
}

func Infof(format string, args ...any) {
	write("INFO", fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...any) {
	write("ERROR", fmt.Sprintf(format, args...))
}

func write(level, msg string) {
	mu.Lock()
	path := logPath
	mu.Unlock()

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()

	line := fmt.Sprintf("%s [%s] %s\n", time.Now().Format(time.RFC3339), level, msg)
	_, _ = f.WriteString(line)
}
