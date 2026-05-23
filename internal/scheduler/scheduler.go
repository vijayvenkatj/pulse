package scheduler

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/vijayvenkatj/pulse/internal/aggregator"
	"github.com/vijayvenkatj/pulse/internal/core"
	"github.com/vijayvenkatj/pulse/internal/pacer"
)

func Scheduler(ctx context.Context, httpClient *http.Client, config core.Config) []core.Result {

	if config.Duration > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, config.Duration)
		defer cancel()
	}

	// Job channel
	jobChan := make(chan core.Job, config.Requests)

	// Result channel
	resultChan := make(chan core.Result, config.Requests)

	var results []core.Result
	done := make(chan struct{})
	go func() {
		results = aggregator.Aggregator(resultChan)
		close(done)
	}()

	// Create Workers
	var wg sync.WaitGroup
	for _ = range config.Concurrency {
		wg.Go(func() {
			Worker(ctx, jobChan, resultChan, httpClient)
		})
	}

	// Create Pacer
	jobFn := core.CreateJobFactory("http://localhost:9001", "GET", []byte{}, nil)
	interval := int(time.Second) / config.RPS
	go func() {
		pacer.Pacer(ctx, time.Duration(interval), config.Requests, jobFn, jobChan)
		close(jobChan)
	}()

	wg.Wait()
	close(resultChan)
	<-done

	return results
}
