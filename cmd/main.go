package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/vijayvenkatj/pulse/internal/aggregator"
	"github.com/vijayvenkatj/pulse/internal/core"
	"github.com/vijayvenkatj/pulse/internal/scheduler"
)

func main() {
	config := core.Config{}
	var durationStr string
	var payloadStr string

	flag.IntVar(&config.Concurrency, "c", 10, "Number of concurrent workers")
	flag.IntVar(&config.Requests, "n", 1000, "Total number of requests")
	flag.IntVar(&config.RPS, "r", 200, "Requests per second")
	flag.StringVar(&durationStr, "d", "5s", "Duration of the test (e.g. 10s, 1m)")

	flag.StringVar(&config.Method, "m", "GET", "HTTP method")
	flag.StringVar(&payloadStr, "p", "", "Request payload string")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: pulse [options] <url>\n\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Error: URL is required")
		flag.Usage()
		os.Exit(1)
	}

	config.URL = flag.Arg(0)
	config.Payload = []byte(payloadStr)
	config.Duration, _ = time.ParseDuration(durationStr)

	results := scheduler.Scheduler(context.Background(), &http.Client{}, config)
	fmt.Print("\n" + aggregator.Report(results))
}
