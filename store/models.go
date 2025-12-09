package store

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// HTTPCheck represents a single HTTP check result
type HTTPCheck struct {
	// Basic endpoint info
	ID             string        `db:"id"`
	Name           string        `db:"name"`
	URL            string        `db:"url"`
	Method         string        `db:"method"`
	Interval       time.Duration `db:"interval_seconds"`
	Timeout        time.Duration `db:"timeout_seconds"`
	ExpectedStatus int           `db:"expected_status"`

	// Additional config
	Headers         map[string]string `db:"headers"`
	Body            string            `db:"body"`
	MustMatchStatus bool              `db:"must_match_status"`
	BodyContains    string            `db:"body_contains"`
	BodyRegex       string            `db:"body_regex"`

	// Results
	Status          string            `db:"status"` // "up", "degraded", "down", "unreachable"
	StatusCode      int               `db:"status_code"`
	ContentLength   int               `db:"content_length"`
	ResponseHeaders map[string]string `db:"response_headers"`
	ResponseBody    string            `db:"response_body"`
	ErrorMessage    string            `db:"error_message"`
	Duration        time.Duration     `db:"duration_ms"`

	// Timing metrics
	DNSDuration       time.Duration `db:"dns_duration_ms"`
	TLSDuration       time.Duration `db:"tls_duration_ms"`
	ConnectDuration   time.Duration `db:"connect_duration_ms"`
	FirstByteDuration time.Duration `db:"first_byte_duration_ms"`
	DownloadDuration  time.Duration `db:"download_duration_ms"`
	SSLValid          bool          `db:"ssl_valid"`
	SSLExpiryDays     int           `db:"ssl_expiry_days"`

	// Timestamps
	CheckedAt time.Time `db:"checked_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// For database operations (uses seconds/milliseconds)
type DBHTTPCheck struct {
	ID              string `db:"id"`
	Name            string `db:"name"`
	URL             string `db:"url"`
	Method          string `db:"method"`
	IntervalSeconds int64  `db:"interval_seconds"`
	TimeoutSeconds  int64  `db:"timeout_seconds"`
	ExpectedStatus  int    `db:"expected_status"`

	Headers         JSONMap `db:"headers"`
	Body            string  `db:"body"`
	MustMatchStatus bool    `db:"must_match_status"`
	BodyContains    string  `db:"body_contains"`
	BodyRegex       string  `db:"body_regex"`

	Status          string  `db:"status"`
	StatusCode      int     `db:"status_code"`
	ContentLength   int     `db:"content_length"`
	ResponseHeaders JSONMap `db:"response_headers"`
	ResponseBody    string  `db:"response_body"`
	ErrorMessage    string  `db:"error_message"`
	DurationMS      int64   `db:"duration_ms"`

	DNSDurationMS       int64 `db:"dns_duration_ms"`
	TLSDurationMS       int64 `db:"tls_duration_ms"`
	ConnectDurationMS   int64 `db:"connect_duration_ms"`
	FirstByteDurationMS int64 `db:"first_byte_duration_ms"`
	DownloadDurationMS  int64 `db:"download_duration_ms"`
	SSLValid            bool  `db:"ssl_valid"`
	SSLExpiryDays       int   `db:"ssl_expiry_days"`

	CheckedAt time.Time `db:"checked_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Convert to/from DB model
func (h *HTTPCheck) ToDB() *DBHTTPCheck {
	return &DBHTTPCheck{
		ID:                  h.ID,
		Name:                h.Name,
		URL:                 h.URL,
		Method:              h.Method,
		IntervalSeconds:     int64(h.Interval.Seconds()),
		TimeoutSeconds:      int64(h.Timeout.Seconds()),
		ExpectedStatus:      h.ExpectedStatus,
		Headers:             JSONMap(h.Headers),
		Body:                h.Body,
		MustMatchStatus:     h.MustMatchStatus,
		BodyContains:        h.BodyContains,
		BodyRegex:           h.BodyRegex,
		Status:              h.Status,
		StatusCode:          h.StatusCode,
		ContentLength:       h.ContentLength,
		ResponseHeaders:     JSONMap(h.ResponseHeaders),
		ResponseBody:        h.ResponseBody,
		ErrorMessage:        h.ErrorMessage,
		DurationMS:          h.Duration.Milliseconds(),
		DNSDurationMS:       h.DNSDuration.Milliseconds(),
		TLSDurationMS:       h.TLSDuration.Milliseconds(),
		ConnectDurationMS:   h.ConnectDuration.Milliseconds(),
		FirstByteDurationMS: h.FirstByteDuration.Milliseconds(),
		DownloadDurationMS:  h.DownloadDuration.Milliseconds(),
		SSLValid:            h.SSLValid,
		SSLExpiryDays:       h.SSLExpiryDays,
		CheckedAt:           h.CheckedAt,
		CreatedAt:           h.CreatedAt,
		UpdatedAt:           h.UpdatedAt,
	}
}

func (db *DBHTTPCheck) ToDomain() *HTTPCheck {
	return &HTTPCheck{
		ID:                db.ID,
		Name:              db.Name,
		URL:               db.URL,
		Method:            db.Method,
		Interval:          time.Duration(db.IntervalSeconds) * time.Second,
		Timeout:           time.Duration(db.TimeoutSeconds) * time.Second,
		ExpectedStatus:    db.ExpectedStatus,
		Headers:           map[string]string(db.Headers),
		Body:              db.Body,
		MustMatchStatus:   db.MustMatchStatus,
		BodyContains:      db.BodyContains,
		BodyRegex:         db.BodyRegex,
		Status:            db.Status,
		StatusCode:        db.StatusCode,
		ContentLength:     db.ContentLength,
		ResponseHeaders:   map[string]string(db.ResponseHeaders),
		ResponseBody:      db.ResponseBody,
		ErrorMessage:      db.ErrorMessage,
		Duration:          time.Duration(db.DurationMS) * time.Millisecond,
		DNSDuration:       time.Duration(db.DNSDurationMS) * time.Millisecond,
		TLSDuration:       time.Duration(db.TLSDurationMS) * time.Millisecond,
		ConnectDuration:   time.Duration(db.ConnectDurationMS) * time.Millisecond,
		FirstByteDuration: time.Duration(db.FirstByteDurationMS) * time.Millisecond,
		DownloadDuration:  time.Duration(db.DownloadDurationMS) * time.Millisecond,
		SSLValid:          db.SSLValid,
		SSLExpiryDays:     db.SSLExpiryDays,
		CheckedAt:         db.CheckedAt,
		CreatedAt:         db.CreatedAt,
		UpdatedAt:         db.UpdatedAt,
	}
}

// JSONMap for JSONB fields
type JSONMap map[string]string

func (jm JSONMap) Value() (driver.Value, error) {
	if len(jm) == 0 {
		return "{}", nil
	}
	return json.Marshal(jm)
}

func (jm *JSONMap) Scan(value any) error {
	if value == nil {
		*jm = JSONMap{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unsupported type for JSONMap: %T", value)
	}
	return json.Unmarshal(b, jm)
}
