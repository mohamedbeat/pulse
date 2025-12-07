package httpchecker

import (
	"context"
	"fmt"
	"github.com/mohamedbeat/pulse/common"
	"net/http"
	"time"
)

// HTTPClient defines the interface for HTTP operations
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// DefaultHTTPClient implements HTTPClient using the real http.Client
type DefaultHTTPClient struct {
	client *http.Client
}

type HTTPChecker struct {
	client HTTPClient
}

func NewDefaultHTTPClient() *DefaultHTTPClient {
	return &DefaultHTTPClient{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

func (d *DefaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return d.client.Do(req)
}

// NewHTTPChecker creates a new checker with default HTTP client
func NewHTTPChecker() *HTTPChecker {
	return &HTTPChecker{
		client: NewDefaultHTTPClient(),
	}
}

// NewHTTPCheckerWithClient allows injection of custom HTTP client (for testing)
func NewHTTPCheckerWithClient(client HTTPClient) *HTTPChecker {
	return &HTTPChecker{
		client: client,
	}
}

func (h *HTTPChecker) Check(ctx context.Context, endpoint common.Endpoint) common.Result {
	start := time.Now()

	ctx, cancel := context.WithTimeout(ctx, endpoint.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, endpoint.Method, endpoint.URL, nil)

	result := common.Result{
		URL:      endpoint.URL,
		Messages: make([]string, 0),
	}

	if err != nil {
		result.Status = common.StatusUnreachable
		result.Error = err.Error()
		result.Timestamp = time.Now()

		return result
	}

	// Add custom headers
	if len(endpoint.Headers) > 0 {
		// Debug("Setting headers", "endpoint", endpoint.URL, "headers", endpoint.Headers)
		for k, v := range endpoint.Headers {
			req.Header.Set(k, v)
		}
	}

	// Starting the request
	resp, err := h.client.Do(req)
	elapsed := time.Since(start)

	result.Elapsed = int(elapsed.Milliseconds())
	result.Timestamp = time.Now()

	if err != nil {
		result.Status = common.StatusUnreachable
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.Status = common.StatusUp
	} else if resp.StatusCode >= 500 {
		result.Status = common.StatusDown
	} else {
		result.Status = common.StatusDegraded
	}
	result.StatusCode = resp.StatusCode

	// Check if resp.StatusCode must match
	if endpoint.MustMatchStatus && resp.StatusCode != endpoint.ExpectedStatus {
		fmt.Println("unexpected status detected")
		// Mark as degraded if status code doesn't match expected strict status.
		if result.Status == common.StatusUp {
			result.Status = common.StatusDegraded
		}
		result.Messages = append(result.Messages, common.UnexpectedStatusCodeMessage)
	}

	if endpoint.MaxLatency > 0 && elapsed > endpoint.MaxLatency {
		fmt.Println("latency detected")
		// If we were "up" purely by status code, treat high latency as degraded.
		if result.Status == common.StatusUp {
			result.Status = common.StatusDegraded
		}
		result.Messages = append(result.Messages, common.UnexpectedLatencyMessage)
	}
	return result
}
