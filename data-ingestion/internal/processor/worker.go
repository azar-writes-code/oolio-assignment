package processor

import (
	"context"
	"log"
	"sync"

	"github.com/azar-writes-code/oolio-data-ingestion/internal/storage"
)

type Ingester struct {
	DB        *storage.BadgerClient
	BatchSize int
}

func (n *Ingester) Worker(ctx context.Context, jobs <-chan string, bitOffset int, wg *sync.WaitGroup) {
	defer wg.Done()

	batch := make([]string, 0, n.BatchSize)

	for code := range jobs {
		batch = append(batch, code)

		if len(batch) >= n.BatchSize {
			if err := n.DB.MarkBatch(batch, bitOffset); err != nil {
				log.Printf("Error writing batch: %v", err)
			}
			batch = batch[:0]
		}
	}
	if len(batch) > 0 {
		if err := n.DB.MarkBatch(batch, bitOffset); err != nil {
			log.Printf("Error writing final batch: %v", err)
		}
	}
}
