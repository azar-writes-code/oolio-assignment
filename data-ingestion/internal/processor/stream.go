package processor

import (
	"bufio"
	"context"
	"io"
	"sync"

	"github.com/cespare/xxhash/v2"
)

func (n *Ingester) Process(ctx context.Context, r io.Reader, bitOffset int) error {
	scanner := bufio.NewScanner(r)

	const numOfWorkers = 16
	jobChannels := make([]chan string, numOfWorkers)
	var wg sync.WaitGroup

	for w := 0; w < numOfWorkers; w++ {
		jobChannels[w] = make(chan string, 5000)
		wg.Add(1)
		go n.Worker(ctx, jobChannels[w], bitOffset, &wg)
	}

	for scanner.Scan() {
		code := scanner.Text()
		if len(code) >= 8 && len(code) <= 10 {
			// Hash the code to find a deterministic worker ID
			workerID := xxhash.Sum64String(code) % numOfWorkers
			jobChannels[workerID] <- code
		}
	}

	// Close all channels to signal workers to finish
	for w := 0; w < numOfWorkers; w++ {
		close(jobChannels[w])
	}
	
	wg.Wait()
	return scanner.Err()
}