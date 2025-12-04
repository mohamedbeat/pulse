package main

import "time"

const (
	MethodGet   = "GET"
	MethodPost  = "POST"
	MethodPatch = "PATCH"
	MethodPut   = "PUT"
)

type Endpoint struct {
	Name            string
	URL             string
	Method          string
	Timeout         time.Duration
	Interval        time.Duration
	Headers         map[string]string
	ExpectedStatus  int
	MustMatchStatus bool
	Type            string        // http, tcp, dns ..
	BodyContains    string        `yaml:"body_contains"`
	BodyRegex       string        `yaml:"body_regex"`
	MaxLatency      time.Duration `yaml:"max_latency"`
}

type Result struct {
	URL        string
	Status     string // "up", "down", "degraded"
	StatusCode int
	// Latency      time.Duration
	Timestamp    time.Time
	ResponseTime int // milliseconds
	Error        string
	Message      string `json:"message,omitempty"`
}

// type Result struct {
//     EndpointID string        `json:"endpoint_id"`
//     Timestamp  time.Time     `json:"timestamp"`
//     Status     string        `json:"status"` // "healthy", "unhealthy", "timeout"
//     Latency    time.Duration `json:"latency"`
//     StatusCode int           `json:"status_code,omitempty"`
//     Error      string        `json:"error,omitempty"`
//     Message    string        `json:"message,omitempty"`
// }
