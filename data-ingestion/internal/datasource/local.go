package datasource

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type LocalProvider struct {
	Dir string
}

func (p *LocalProvider) ListSources(ctx context.Context) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(p.Dir, "*"))
	return files, err
}

func (p *LocalProvider) GetStream(ctx context.Context, sourceID string) (io.ReadCloser, error) {
	return os.Open(sourceID)
}