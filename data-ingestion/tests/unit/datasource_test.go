package unit

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/azar-writes-code/oolio-data-ingestion/internal/datasource"
	"github.com/stretchr/testify/assert"
)

func TestLocalProvider(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "datasource-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a dummy file
	filename := filepath.Join(tempDir, "source1.txt")
	content := "content1\ncontent2"
	err = os.WriteFile(filename, []byte(content), 0644)
	assert.NoError(t, err)

	provider := &datasource.LocalProvider{Dir: tempDir}
	ctx := context.Background()

	// Test ListSources
	sources, err := provider.ListSources(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(sources))
	assert.Equal(t, filename, sources[0])

	// Test GetStream
	stream, err := provider.GetStream(ctx, filename)
	assert.NoError(t, err)
	defer stream.Close()

	readContent, err := io.ReadAll(stream)
	assert.NoError(t, err)
	assert.Equal(t, content, string(readContent))
}
