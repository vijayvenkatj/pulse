package scheduler

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/vijayvenkatj/pulse/internal/core"
)

// Worker takes in the jobs and execs them and sends the results to the aggregator through the resultChan
func Worker(ctx context.Context, jobChan <-chan core.Job, resultChan chan<- core.Result, httpClient *http.Client) {
	for {
		select {
		case <-ctx.Done():
			return

		case job, ok := <-jobChan:
			if !ok {
				return
			}

			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			// Creating a new request with our client body and headers
			reader := bytes.NewReader(job.Body)
			req, err := http.NewRequestWithContext(ctx, job.Method, job.URL.String(), reader)
			if err != nil {
				log.Printf("ERROR: job %d skipped, %v", job.ID, err)
				continue
			}
			req.Header = job.Headers.Clone()

			result := core.Result{
				JobID:     job.ID,
				TimeStamp: time.Now(),
				BytesOut:  int64(len(job.Body)),
			}

			// This is to measure the metrics of the request
			start := time.Now()
			resp, reqErr := httpClient.Do(req)
			result.TTFB = time.Since(start)
			if reqErr != nil {
				result.Err = reqErr
				select {
				case <-ctx.Done():
					return
				case resultChan <- result:

				}
				continue
			}

			bytesIn, bodyErr := io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			if bodyErr != nil {
				result.Err = bodyErr
			}

			result.Latency = time.Since(start)
			result.BytesIn = bytesIn
			result.StatusCode = resp.StatusCode

			// Sending the result through the chan and checking for ctx cancellation in case.
			select {
			case <-ctx.Done():
				return
			case resultChan <- result:

			}
		}
	}
}
