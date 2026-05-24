package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	_ "net/http/pprof"

	"github.com/vijayvenkatj/pulse/internal/aggregator"
	"github.com/vijayvenkatj/pulse/internal/core"
	"github.com/vijayvenkatj/pulse/internal/scheduler"
)

func main() {

	runtime.SetBlockProfileRate(1)
	go func() {
		fmt.Println(http.ListenAndServe(
			"localhost:6060",
			nil,
		))
	}()

	config := core.Config{}
	var durationStr string
	var payloadStr string

	flag.IntVar(&config.Concurrency, "c", 10, "Number of concurrent workers")
	flag.IntVar(&config.Requests, "n", 0, "Total number of requests")
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

	metrics := scheduler.Scheduler(context.Background(), &http.Client{}, config)
	fmt.Print("\n" + aggregator.Report(metrics))
}
