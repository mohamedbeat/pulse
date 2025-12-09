package repository

import "time"

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

type HTTPCheckStats struct {
	TotalEndpoints   int       `db:"total_endpoints"`
	TotalChecks      int       `db:"total_checks"`
	SuccessfulChecks int       `db:"successful_checks"`
	FailedChecks     int       `db:"failed_checks"`
	DegradedChecks   int       `db:"degraded_checks"`
	AvgLatencyMS     float64   `db:"avg_latency_ms"`
	LastCheckedAt    time.Time `db:"last_checked_at"`
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
