package faketime

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Provider interface {
	Get(ctx context.Context) (string, error)
}

type FileProvider struct {
	filePath string
}

var _ Provider = (*FileProvider)(nil)

func NewFileProvider(filePath string) *FileProvider {
	return &FileProvider{filePath: filePath}
}

var (
	ErrFileRead = fmt.Errorf("failed to read faketime file")
)

func (p *FileProvider) Get(ctx context.Context) (string, error) {
	stat, err := os.Stat(p.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			slog.DebugContext(ctx, "time file does not exist, proceeding without setting time", "file", p.filePath)
			return "", nil
		} else {
			return "", fmt.Errorf("%w: %v", ErrFileRead, err)
		}
	}
	if stat.IsDir() {
		return "", fmt.Errorf("%w: path is a directory, not a file: %s", ErrFileRead, p.filePath)
	}

	data, err := os.ReadFile(p.filePath)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrFileRead, err)
	}
	return strings.TrimSpace(string(data)), nil
}
