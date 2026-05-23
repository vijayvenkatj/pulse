package core

import "time"

type Config struct {
	Concurrency int
	Requests    int
	RPS         int
	Duration    time.Duration
	Method      string
	URL         string
	Payload     []byte
}
