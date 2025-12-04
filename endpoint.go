package main

import "time"

const (
	MethodGet   = "GET"
	MethodPost  = "POST"
	MethodPatch = "PATCH"
	MethodPut   = "PUT"
)

type Endpoint struct {
	Name            string            `json:"name" yaml:"name"`
	URL             string            `json:"url" yaml:"url"`
	Method          string            `json:"method" yaml:"method"`
	Timeout         time.Duration     `json:"timeout" yaml:"timeout"`
	Interval        time.Duration     `json:"interval" yaml:"interval"`
	Headers         map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	ExpectedStatus  int               `json:"expected_status" yaml:"expected_status"`
	MustMatchStatus bool              `json:"must_match_status" yaml:"must_match_status"`
	Type            string            `json:"type" yaml:"type"` // http, tcp, dns ..
	BodyContains    string            `json:"body_contains,omitempty" yaml:"body_contains,omitempty"`
	BodyRegex       string            `json:"body_regex,omitempty" yaml:"body_regex,omitempty"`
	MaxLatency      time.Duration     `json:"max_latency" yaml:"max_latency"`
}

// Result.Status
const (
	StatusUp          = "up"
	StatusDown        = "down"
	StatusUnreachable = "unreachable"
	StatusDegraded    = "degraded"
)

// Result.Message
const (
	UnexpectedStatusCodeMessage = "UnexpectedStatusCode"
	UnexpectedBodyMessage       = "UnexpectedBody"
	UnexpectedLatencyMessage    = "UnexpectedLatency"
)

type Result struct {
	URL        string    `json:"url" yaml:"url"`
	Status     string    `json:"status" yaml:"status"` // "up", "degraded", "down", "unreachable"
	StatusCode int       `json:"status_code" yaml:"status_code"`
	Timestamp  time.Time `json:"timestamp" yaml:"timestamp"`
	Elapsed    int       `json:"elapsed_ms" yaml:"elapsed_ms"` // milliseconds
	Error      string    `json:"error,omitempty" yaml:"error,omitempty"`
	Message    string    `json:"message,omitempty" yaml:"message,omitempty"`
}
