package aggregator

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/vijayvenkatj/pulse/internal/core"
)

// Report formats a minimal report from the results slice.
func Report(results []core.Result) string {
	var b strings.Builder
	b.WriteString("Pulse Report\n")

	total := len(results)
	if total == 0 {
		b.WriteString("No results to report.\n")
		return b.String()
	}

	var (
		successCount int
		errorCount   int
		bytesIn      int64
		bytesOut     int64
		statusCounts = map[int]int{}
		latencies    []time.Duration
		ttfbs        []time.Duration
		startTime    time.Time
		endTime      time.Time
	)

	for _, result := range results {
		bytesIn += result.BytesIn
		bytesOut += result.BytesOut

		if result.Err != nil {
			errorCount++
		} else {
			successCount++
			latencies = append(latencies, result.Latency)
			ttfbs = append(ttfbs, result.TTFB)
		}

		if result.StatusCode > 0 {
			statusCounts[result.StatusCode]++
		}

		if startTime.IsZero() || result.TimeStamp.Before(startTime) {
			startTime = result.TimeStamp
		}
		endCandidate := result.TimeStamp.Add(result.Latency)
		if endCandidate.After(endTime) {
			endTime = endCandidate
		}
	}

	var duration time.Duration
	if !startTime.IsZero() && endTime.After(startTime) {
		duration = endTime.Sub(startTime)
	}

	b.WriteString("\n")
	tw := tabwriter.NewWriter(&b, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, "Summary\tValue")
	fmt.Fprintf(tw, "Requests\t%d\n", total)
	fmt.Fprintf(tw, "Success\t%d (%.2f%%)\n", successCount, percent(successCount, total))
	fmt.Fprintf(tw, "Errors\t%d\n", errorCount)
	if duration > 0 {
		fmt.Fprintf(tw, "Test duration\t%s\n", formatDuration(duration))
	}
	fmt.Fprintf(tw, "Bytes in\t%s\n", formatBytes(bytesIn))
	fmt.Fprintf(tw, "Bytes out\t%s\n", formatBytes(bytesOut))
	tw.Flush()

	if successCount > 0 {
		b.WriteString("\nLatency (success)\n")
		writeDurationStats(&b, latencies)

		b.WriteString("\nTTFB (success)\n")
		writeDurationStats(&b, ttfbs)
	} else {
		b.WriteString("\nLatency (success)\nNo successful requests.\n")
		b.WriteString("\nTTFB (success)\nNo successful requests.\n")
	}

	b.WriteString("\nStatus codes\n")
	if len(statusCounts) == 0 {
		b.WriteString("No status codes captured.\n")
		return b.String()
	}

	codes := make([]int, 0, len(statusCounts))
	for code := range statusCounts {
		codes = append(codes, code)
	}
	sort.Ints(codes)

	tw = tabwriter.NewWriter(&b, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, "Code\tCount\tPercent")
	for _, code := range codes {
		count := statusCounts[code]
		fmt.Fprintf(tw, "%d\t%d\t%.2f%%\n", code, count, percent(count, total))
	}
	tw.Flush()

	return b.String()
}

func writeDurationStats(b *strings.Builder, values []time.Duration) {
	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
	tw := tabwriter.NewWriter(b, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, "min\tavg\tp50\tp90\tp99\tmax")
	fmt.Fprintf(
		tw,
		"%s\t%s\t%s\t%s\t%s\t%s\n",
		formatDuration(values[0]),
		formatDuration(avgDuration(values)),
		formatDuration(percentile(values, 50)),
		formatDuration(percentile(values, 90)),
		formatDuration(percentile(values, 99)),
		formatDuration(values[len(values)-1]),
	)
	tw.Flush()
}

func avgDuration(values []time.Duration) time.Duration {
	if len(values) == 0 {
		return 0
	}
	var total time.Duration
	for _, value := range values {
		total += value
	}
	return time.Duration(int64(total) / int64(len(values)))
}

func percentile(values []time.Duration, p float64) time.Duration {
	if len(values) == 0 {
		return 0
	}
	if p <= 0 {
		return values[0]
	}
	if p >= 100 {
		return values[len(values)-1]
	}
	index := int(math.Ceil((p / 100) * float64(len(values))))
	if index <= 0 {
		return values[0]
	}
	if index > len(values) {
		return values[len(values)-1]
	}
	return values[index-1]
}

func percent(part, total int) float64 {
	if total == 0 {
		return 0
	}
	return (float64(part) / float64(total)) * 100
}

func formatDuration(value time.Duration) string {
	if value <= 0 {
		return "0s"
	}
	if value < time.Millisecond {
		return value.Round(time.Microsecond).String()
	}
	return value.Round(time.Millisecond).String()
}

func formatBytes(value int64) string {
	if value < 1024 {
		return fmt.Sprintf("%d B", value)
	}
	units := []string{"KB", "MB", "GB", "TB", "PB"}
	size := float64(value)
	unit := 0
	for size >= 1024 && unit < len(units)-1 {
		size /= 1024
		unit++
	}
	return fmt.Sprintf("%.2f %s", size, units[unit])
}
