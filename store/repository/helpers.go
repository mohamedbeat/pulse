package repository

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// For summary queries
type EndpointSummary struct {
	Name       string    `db:"name"`
	URL        string    `db:"url"`
	Status     string    `db:"status"`
	StatusCode int       `db:"status_code"`
	DurationMS int64     `db:"duration_ms"`
	CheckedAt  time.Time `db:"checked_at"`
	Error      string    `db:"error_message"`
}

func (e *EndpointSummary) String() string {
	var sb strings.Builder

	// Get the reflect value of the struct (dereference the pointer)
	v := reflect.ValueOf(e).Elem()
	t := v.Type() // Get type from the value

	// Loop through all fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name

		sb.WriteString(fmt.Sprintf("%s: %v,\n ",
			fieldName,
			field.Interface(),
		))
	}

	return sb.String()
}

type HTTPCheckStats struct {
	TotalEndpoints   int       `db:"total_endpoints"`
	TotalChecks      int       `db:"total_checks"`
	SuccessfulChecks int       `db:"successful_checks"`
	FailedChecks     int       `db:"failed_checks"`
	DegradedChecks   int       `db:"degraded_checks"`
	AvgLatencyMS     float64   `db:"avg_latency_ms"`
	LastCheckedAt    time.Time `db:"last_checked_at"`
}

func (h *HTTPCheckStats) String() string {
	var sb strings.Builder

	// Get the reflect value of the struct (dereference the pointer)
	v := reflect.ValueOf(h).Elem()
	t := v.Type() // Get type from the value

	// Loop through all fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name

		sb.WriteString(fmt.Sprintf("%s: %v,\n ",
			fieldName,
			field.Interface(),
		))
	}

	return sb.String()
}

type CheckFilters struct {
	Name      string
	Status    []string
	StartTime time.Time
	EndTime   time.Time
	Limit     int
	Offset    int
}

type FailingChecksOptions struct {
	Since        time.Time
	Limit        int
	EndpointName string
}
