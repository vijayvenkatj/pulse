# Pulse

Pulse is a small HTTP load generator that sends requests at a target RPS, captures latency/TTFB, and prints a compact report.

## Requirements

- Go 1.26+

## Install

```bash
git clone https://github.com/vijayvenkatj/pulse.git
cd pulse
```

## Usage

```bash
go run ./cmd -n 50 -c 5 -r 20 -d 5s http://127.0.0.1:8080/ok
```

### Options

| Flag | Description | Default |
| --- | --- | --- |
| `-c` | Concurrent workers | `10` |
| `-n` | Total requests | `1000` |
| `-r` | Requests per second | `200` |
| `-d` | Test duration (e.g. `10s`, `1m`) | `5s` |
| `-m` | HTTP method | `GET` |
| `-p` | Request body string | `""` |

## Example output

```
Pulse Report

Summary        Value
Requests       50
Success        50 (100.00%)
Errors         0
Test duration  2.45s
Bytes in       1.51 MB
Bytes out      0 B

Latency (success)
min    avg    p50    p90  p99  max
441µs  827µs  778µs  1ms  1ms  1ms

TTFB (success)
min    avg    p50    p90  p99  max
432µs  804µs  748µs  1ms  1ms  1ms

Status codes
Code  Count  Percent
200   50     100.00%
```
