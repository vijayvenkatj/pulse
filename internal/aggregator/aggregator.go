package aggregator

import (
	"github.com/vijayvenkatj/pulse/internal/core"
)

// Aggregator takes in results and makes a Cummulated Result Slice.
func Aggregator(resultChan <-chan core.Result, results []core.Result) {
	i := 0
	for result := range resultChan {
		results[i] = result
		i++
	}
}
