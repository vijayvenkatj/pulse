package server

import "time"

type Config struct {
	Concurrency int
	Requests    int
	RPS         int
	Duration    time.Duration
}
