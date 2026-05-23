package aggregator

import (
	"github.com/vijayvenkatj/pulse/internal/core"
)

// Aggregator takes in results and makes a Cummulated Result Slice.
func Aggregator(resultChan <-chan core.Result) []core.Result {

	results := []core.Result{}

	for result := range resultChan {
		results = append(results, result)
	}

	return results
}
