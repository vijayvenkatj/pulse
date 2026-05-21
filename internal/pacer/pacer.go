package pacer

import (
	"context"
	"time"

	"github.com/vijayvenkatj/pulse/internal/core"
)

// Pacer is responsible for generating the Jobs WRT the Clients requirements.
func Pacer(ctx context.Context, interval time.Duration, jobFn func() core.Job, jobChan chan<- core.Job) {

	// Makes a Job once per interval, So for X RPS, we get X jobs per sec.
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

			// JobFn() makes a Job with internal metrics.
			job := jobFn()

			// Nested select for context aware job scheduling.
			select {
			case <-ctx.Done():
				return
			case jobChan <- job:
			}
		}
	}
}
