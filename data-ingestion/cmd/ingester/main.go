package main

import (
	"context"
	"log"

	"github.com/azar-writes-code/oolio-data-ingestion/internal/datasource"
	"github.com/azar-writes-code/oolio-data-ingestion/internal/processor"
	"github.com/azar-writes-code/oolio-data-ingestion/internal/storage"
)

func main() {
	ctx := context.Background()
	
	// Badger will create this directory to store the SSD-backed data
	bStore, err := storage.NewBadgerDB("../badger-data")
	if err != nil {
		log.Fatalf("Failed to initialize BadgerDB: %v", err)
	}
	defer bStore.Close()

	provider := &datasource.LocalProvider{Dir: "data"}
	
	// Batching 5000 records at a time is highly efficient for Badger
	ingester := &processor.Ingester{DB: bStore, BatchSize: 5000}

	files, err := provider.ListSources(ctx)
	if err != nil {
		log.Fatalf("Failed to list sources: %v", err)
	}
	if len(files) == 0 {
		log.Println("No source files found in the directory.")
	}

	for i, file := range files {
		log.Printf("Processing %s with bit offset %d", file, i)
		stream, err := provider.GetStream(ctx, file)
		if err != nil {
			log.Printf("Failed to get stream for %s: %v", file, err)
			continue
		}
		
		if err := ingester.Process(ctx, stream, i); err != nil {
			log.Printf("Error processing %s: %v", file, err)
		}
		stream.Close()
	}
	log.Println("Ingestion Complete. Data is safely stored in ./badger-data")
}