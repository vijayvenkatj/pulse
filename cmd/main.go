package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/vijayvenkatj/pulse/internal/scheduler"
	"github.com/vijayvenkatj/pulse/internal/server"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config := server.Config{
		Concurrency: 30,
		Requests:    2000,
		RPS:         1000,
		Duration:    6 * time.Second,
	}

	httpClient := http.Client{}

	start := time.Now()
	scheduler.Scheduler(ctx, &httpClient, config)
	end := time.Since(start)

	log.Println("TIME: ", end)
}
