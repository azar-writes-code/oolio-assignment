package datasource

import (
	"context"
	"io"
)

type Provider interface {
	GetStream(ctx context.Context, sourceID string) (io.ReadCloser, error)
	ListSources(ctx context.Context) ([]string, error)
}