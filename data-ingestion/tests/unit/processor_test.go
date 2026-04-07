package unit

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/azar-writes-code/oolio-data-ingestion/internal/processor"
	"github.com/stretchr/testify/assert"
)

type mockStorage struct {
	batches [][]string
	offsets []int
	mu      sync.Mutex
}

func (m *mockStorage) MarkBatch(codes []string, bitOffset int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	batchCopy := make([]string, len(codes))
	copy(batchCopy, codes)
	m.batches = append(m.batches, batchCopy)
	m.offsets = append(m.offsets, bitOffset)
	return nil
}

func TestIngester_Process(t *testing.T) {
	storage := &mockStorage{}
	ingester := &processor.Ingester{
		DB:        storage,
		BatchSize: 2,
	}

	input := `abc
	valid123
	toolongcodes
	valid456
	short
	valid789
	`
	ctx := context.Background()
	reader := strings.NewReader(input)

	err := ingester.Process(ctx, reader, 1)
	assert.NoError(t, err)

	// Codes should be filtered (len >= 8 && len <= 10)
	// valid123 (8), valid456 (8), valid789 (8)
	// 'abc' (3) - too short
	// 'toolongcodes' (12) - too long
	// 'short' (5) - too short

	allCodes := []string{}
	for _, batch := range storage.batches {
		allCodes = append(allCodes, batch...)
	}

	assert.ElementsMatch(t, []string{"valid123", "valid456", "valid789"}, allCodes)
	for _, offset := range storage.offsets {
		assert.Equal(t, 1, offset)
	}
}

func TestIngester_Worker(t *testing.T) {
	storage := &mockStorage{}
	ingester := &processor.Ingester{
		DB:        storage,
		BatchSize: 3,
	}

	jobs := make(chan string, 10)
	var wg sync.WaitGroup
	wg.Add(1)

	ctx := context.Background()
	go ingester.Worker(ctx, jobs, 2, &wg)

	codes := []string{"code1111", "code2222", "code3333", "code4444", "code5555"}
	for _, c := range codes {
		jobs <- c
	}
	close(jobs)
	wg.Wait()

	// With BatchSize=3 and 5 codes:
	// Batch 1: code1111, code2222, code3333
	// Batch 2: code4444, code5555 (final batch)
	assert.Equal(t, 2, len(storage.batches))
	assert.Equal(t, []string{"code1111", "code2222", "code3333"}, storage.batches[0])
	assert.Equal(t, []string{"code4444", "code5555"}, storage.batches[1])
	assert.Equal(t, 2, storage.offsets[0])
	assert.Equal(t, 2, storage.offsets[1])
}
