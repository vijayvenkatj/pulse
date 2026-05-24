package aggregator

import (
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
)

// Report formats a minimal report from rolling metrics.
func Report(metrics *Metrics) string {
	var b strings.Builder
	b.WriteString("Pulse Report\n")

	if metrics == nil || metrics.total == 0 {
		b.WriteString("No results to report.\n")
		return b.String()
	}

	var duration time.Duration
	if !metrics.startTime.IsZero() && metrics.endTime.After(metrics.startTime) {
		duration = metrics.endTime.Sub(metrics.startTime)
	}

	b.WriteString("\n")
	tw := tabwriter.NewWriter(&b, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, "Summary\tValue")
	fmt.Fprintf(tw, "Requests\t%d\n", metrics.total)
	fmt.Fprintf(tw, "Success\t%d (%.2f%%)\n", metrics.success, percent(metrics.success, metrics.total))
	fmt.Fprintf(tw, "Errors\t%d\n", metrics.errors)
	if duration > 0 {
		fmt.Fprintf(tw, "Test duration\t%s\n", formatDuration(duration))
	}
	fmt.Fprintf(tw, "Bytes in\t%s\n", formatBytes(metrics.bytesIn))
	fmt.Fprintf(tw, "Bytes out\t%s\n", formatBytes(metrics.bytesOut))
	tw.Flush()

	if metrics.success > 0 {
		b.WriteString("\nLatency (success)\n")
		writeHistogramStats(&b, metrics.latencyHist)

		b.WriteString("\nTTFB (success)\n")
		writeHistogramStats(&b, metrics.ttfbHist)
	} else {
		b.WriteString("\nLatency (success)\nNo successful requests.\n")
		b.WriteString("\nTTFB (success)\nNo successful requests.\n")
	}

	b.WriteString("\nStatus codes\n")
	if len(metrics.statusCounts) == 0 {
		b.WriteString("No status codes captured.\n")
		return b.String()
	}

	codes := make([]int, 0, len(metrics.statusCounts))
	for code := range metrics.statusCounts {
		codes = append(codes, code)
	}
	sort.Ints(codes)

	tw = tabwriter.NewWriter(&b, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, "Code\tCount\tPercent")
	for _, code := range codes {
		count := metrics.statusCounts[code]
		fmt.Fprintf(tw, "%d\t%d\t%.2f%%\n", code, count, percent(count, metrics.total))
	}
	tw.Flush()

	return b.String()
}

func writeHistogramStats(b *strings.Builder, hist *hdrhistogram.Histogram) {
	if hist == nil || hist.TotalCount() == 0 {
		b.WriteString("No samples.\n")
		return
	}

	tw := tabwriter.NewWriter(b, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, "min\tavg\tp50\tp90\tp99\tmax")
	fmt.Fprintf(
		tw,
		"%s\t%s\t%s\t%s\t%s\t%s\n",
		formatDuration(time.Duration(hist.Min())*time.Microsecond),
		formatDuration(time.Duration(hist.Mean())*time.Microsecond),
		formatDuration(time.Duration(hist.ValueAtQuantile(50))*time.Microsecond),
		formatDuration(time.Duration(hist.ValueAtQuantile(90))*time.Microsecond),
		formatDuration(time.Duration(hist.ValueAtQuantile(99))*time.Microsecond),
		formatDuration(time.Duration(hist.Max())*time.Microsecond),
	)
	tw.Flush()
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
