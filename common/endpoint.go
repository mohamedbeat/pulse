package common

import (
	"fmt"
	"strings"
	"time"
)

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPatch  = "PATCH"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
)

var ValidMethods = map[string]bool{
	MethodGet:    true,
	MethodPost:   true,
	MethodPatch:  true,
	MethodPut:    true,
	MethodDelete: true,
}

// ValidateMethod checks whether m is a valid HTTP method.
// It returns nil if valid, or an error otherwise.
func ValidateMethod(m string) error {
	if m == "" || !ValidMethods[strings.ToUpper(m)] {
		return fmt.Errorf("invalid HTTP method: %q", m)
	}
	return nil
}

const (
	HTTPType = "HTTP"
)

var ValidTypes = map[string]bool{
	HTTPType: true,
}

// ValidateMethod checks whether Endpoint.Type is a valid HTTP method.
// It returns nil if valid, or an error otherwise.
func ValidateType(ep *Endpoint) error {
	ep.Type = strings.ToUpper(ep.Type)
	if ep.Type == "" || !ValidTypes[ep.Type] {
		return fmt.Errorf("invalid type: %q", ep.Type)
	}
	return nil
}

type Endpoint struct {
	Name            string            `mapstructure:"name" json:"name" yaml:"name"`
	URL             string            `mapstructure:"url" json:"url" yaml:"url"`
	Method          string            `mapstructure:"method" json:"method" yaml:"method"`
	Timeout         time.Duration     `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	Interval        time.Duration     `mapstructure:"interval" json:"interval" yaml:"interval"`
	Headers         map[string]string `mapstructure:"headers" json:"headers,omitempty" yaml:"headers,omitempty"`
	Type            string            `mapstructure:"type" json:"type" yaml:"type"` // http, tcp, dns
	ExpectedStatus  int               `mapstructure:"expected_status" json:"expected_status" yaml:"expected_status"`
	MustMatchStatus bool              `mapstructure:"must_match_status" json:"must_match_status" yaml:"must_match_status"`
	BodyContains    string            `mapstructure:"body_contains" json:"body_contains,omitempty" yaml:"body_contains,omitempty"`
	BodyRegex       string            `mapstructure:"body_regex" json:"body_regex,omitempty" yaml:"body_regex,omitempty"`
	MaxLatency      time.Duration     `mapstructure:"max_latency" json:"max_latency" yaml:"max_latency"`
	Retry           int               `mapstructure:"retry" json:"retry" yaml:"retry"`
	RetryCounter    int               //Retry state counter
	LastResult      *Result
}

type Result struct {
	Endpoint   *Endpoint
	URL        string        `json:"url" yaml:"url"`
	Status     string        `json:"status" yaml:"status"` // "up", "degraded", "down", "unreachable"
	StatusCode int           `json:"status_code" yaml:"status_code"`
	Timestamp  time.Time     `json:"timestamp" yaml:"timestamp"`
	Elapsed    time.Duration `json:"elapsed_ms" yaml:"elapsed_ms"` // milliseconds
	Error      string        `json:"error,omitempty" yaml:"error,omitempty"`
	// Message    string    `json:"message,omitempty" yaml:"message,omitempty"`
	Messages []string `json:"messages,omitempty" yaml:"messages,omitempty"`
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
