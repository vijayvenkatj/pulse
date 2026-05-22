package core

import "time"

type Result struct {
	JobID uint32

	Latency    time.Duration
	StatusCode int
	BytesIn    uint32
	BytesOut   uint32

	Err       error
	TimeStamp time.Time
}
