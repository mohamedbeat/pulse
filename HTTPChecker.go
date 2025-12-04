package main

import (
	"context"
	"fmt"
	"log/slog"
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

func (h *HTTPChecker) Check(ctx context.Context, endpoint Endpoint) (Result, error) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(ctx, endpoint.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, endpoint.Method, endpoint.URL, nil)
	if err != nil {

		return Result{Status: "error", Error: err.Error(), URL: endpoint.URL, Timestamp: time.Now()}, err
	}

	// Add custom headers
	for k, v := range endpoint.Headers {
		req.Header.Set(k, v)
	}

	resp, err := h.client.Do(req)
	responseTime := time.Since(start).Milliseconds()

	result := Result{
		ResponseTime: int(responseTime),
		Timestamp:    time.Now(),
	}

	if err != nil {
		result.Status = "down"
		result.Error = err.Error()
		return result, err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.Status = "up"
	} else if resp.StatusCode >= 500 {
		result.Status = "down"
	} else {
		result.Status = "degraded"
	}
	result.StatusCode = resp.StatusCode
	if endpoint.MustMatchStatus && resp.StatusCode != endpoint.ExpectedStatus {
		lg := fmt.Sprintf("Unexpected status response: expected %d but got %d\n", endpoint.ExpectedStatus, resp.StatusCode)
		slog.Warn(lg)
	}
	if endpoint.MaxLatency < time.Duration(responseTime) {
		lg := fmt.Sprintf("Request exceeded max latency: Request took %dms expected %d\n", responseTime, endpoint.MaxLatency.Milliseconds())
		slog.Warn(lg)

	}
	return result, nil
}
