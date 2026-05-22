package core

import "time"

type Result struct {
	JobID uint32

	Latency    time.Duration
	TTFB       time.Duration
	StatusCode int
	BytesIn    int64
	BytesOut   int64

	Err       error
	TimeStamp time.Time
}
