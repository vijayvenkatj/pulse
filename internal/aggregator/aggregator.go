package aggregator

import "github.com/vijayvenkatj/pulse/internal/core"

// Aggregator consumes results and builds rolling metrics.
func Aggregator(resultChan <-chan core.Result) *Metrics {
	metrics := NewMetrics()
	for result := range resultChan {
		metrics.AddResult(result)
	}
	return metrics
}
