package scheduler

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/vijayvenkatj/pulse/internal/aggregator"
	"github.com/vijayvenkatj/pulse/internal/core"
	"github.com/vijayvenkatj/pulse/internal/pacer"
)

func Scheduler(ctx context.Context, httpClient *http.Client, config core.Config) *aggregator.Metrics {

	if config.Duration > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, config.Duration)
		defer cancel()
	}

	// Job channel
	jobChan := make(chan core.Job, config.Requests)

	// Result channel
	resultChan := make(chan core.Result, config.Requests)

	var metrics *aggregator.Metrics
	done := make(chan struct{})
	go func() {
		metrics = aggregator.Aggregator(resultChan)
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
	jobFn := core.CreateJobFactory(config.URL, config.Method, config.Payload, nil)
	if jobFn == nil {
		errorMetrics := aggregator.NewMetrics()
		errorMetrics.AddResult(core.Result{
			Err:       fmt.Errorf("invalid url: %s", config.URL),
			TimeStamp: time.Now(),
		})
		return errorMetrics
	}
	if config.RPS <= 0 {
		errorMetrics := aggregator.NewMetrics()
		errorMetrics.AddResult(core.Result{
			Err:       fmt.Errorf("invalid RPS: %d", config.RPS),
			TimeStamp: time.Now(),
		})
		return errorMetrics
	}
	interval := time.Second / time.Duration(config.RPS)
	go func() {
		pacer.Pacer(ctx, interval, config.Requests, jobFn, jobChan)
		close(jobChan)
	}()

	wg.Wait()
	close(resultChan)
	<-done

	return metrics
}
