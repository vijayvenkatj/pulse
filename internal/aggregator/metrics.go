package aggregator

import (
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
	"github.com/vijayvenkatj/pulse/internal/core"
)

const (
	minLatencyMicros = int64(1)
	maxLatencyMicros = int64((10 * time.Minute) / time.Microsecond)
	histSigFigs      = 3
)

type Metrics struct {
	total        int
	success      int
	errors       int
	bytesIn      int64
	bytesOut     int64
	statusCounts map[int]int
	startTime    time.Time
	endTime      time.Time
	latencyHist  *hdrhistogram.Histogram
	ttfbHist     *hdrhistogram.Histogram
}

func NewMetrics() *Metrics {
	return &Metrics{
		statusCounts: map[int]int{},
		latencyHist:  hdrhistogram.New(minLatencyMicros, maxLatencyMicros, histSigFigs),
		ttfbHist:     hdrhistogram.New(minLatencyMicros, maxLatencyMicros, histSigFigs),
	}
}

func (m *Metrics) AddResult(result core.Result) {
	m.total++
	m.bytesIn += result.BytesIn
	m.bytesOut += result.BytesOut

	if result.Err != nil {
		m.errors++
	} else {
		m.success++
		recordDuration(m.latencyHist, result.Latency)
		recordDuration(m.ttfbHist, result.TTFB)
	}

	if result.StatusCode > 0 {
		m.statusCounts[result.StatusCode]++
	}

	if m.startTime.IsZero() || result.TimeStamp.Before(m.startTime) {
		m.startTime = result.TimeStamp
	}
	endCandidate := result.TimeStamp.Add(result.Latency)
	if endCandidate.After(m.endTime) {
		m.endTime = endCandidate
	}
}

func recordDuration(hist *hdrhistogram.Histogram, value time.Duration) {
	micros := value.Microseconds()
	if micros < minLatencyMicros {
		micros = minLatencyMicros
	}
	if micros > maxLatencyMicros {
		micros = maxLatencyMicros
	}
	_ = hist.RecordValue(micros)
}
