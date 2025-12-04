package main

import (
	"context"
	"net/http"
	"time"
)

type HTTPChecker struct {
	client *http.Client
}

func NewHTTPChecker() *HTTPChecker {
	return &HTTPChecker{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

func (h *HTTPChecker) Check(ctx context.Context, endpoint Endpoint) Result {
	start := time.Now()

	ctx, cancel := context.WithTimeout(ctx, endpoint.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, endpoint.Method, endpoint.URL, nil)

	result := Result{
		URL: endpoint.URL,
	}

	if err != nil {
		result.Status = StatusUnreachable
		result.Error = err.Error()
		result.Timestamp = time.Now()

		return result
	}

	// Add custom headers
	for k, v := range endpoint.Headers {
		req.Header.Set(k, v)
	}

	resp, err := h.client.Do(req)
	elapsed := time.Since(start)

	result.Elapsed = int(elapsed.Milliseconds())
	result.Timestamp = time.Now()

	if err != nil {
		result.Status = StatusUnreachable
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.Status = StatusUp
	} else if resp.StatusCode >= 500 {
		result.Status = StatusDown
	} else {
		result.Status = StatusDegraded
	}
	result.StatusCode = resp.StatusCode

	// Check if resp.StatusCode must match
	if endpoint.MustMatchStatus && resp.StatusCode != endpoint.ExpectedStatus {
		// Mark as degraded if status code doesn't match expected strict status.
		if result.Status == StatusUp {
			result.Status = StatusDegraded
			result.Message = UnexpectedStatusCodeMessage
		}
	}

	if endpoint.MaxLatency > 0 && elapsed > endpoint.MaxLatency {
		// If we were "up" purely by status code, treat high latency as degraded.
		if result.Status == StatusUp {
			result.Status = StatusDegraded
			result.Message = UnexpectedLatencyMessage
		}
	}
	return result
}
