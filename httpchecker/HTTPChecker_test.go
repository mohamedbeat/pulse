package httpchecker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/mohamedbeat/pulse/common"
)

func TestHTTPChecker_Check_Success(t *testing.T) {
	// Create mock client that returns successful response
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify request properties
			// if req.URL.String() != "http://example.com" {
			// 	t.Errorf("Expected URL: http://example.com, got: %s", req.URL.String())
			// }

			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     common.StatusUp,
				Body:       io.NopCloser(strings.NewReader("OK")),
			}, nil
		},
	}

	checker := NewHTTPCheckerWithClient(mockClient)
	endpoint := common.Endpoint{
		URL:     "http://example.com",
		Method:  "GET",
		Timeout: 5 * time.Second,
	}

	result := checker.Check(context.Background(), endpoint)

	if result.Status != common.StatusUp {
		t.Errorf("Expected status Up, got %s", result.Status)
	}

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", result.StatusCode)
	}

	if result.Error != "" {
		t.Errorf("Expected no error, got: %s", result.Error)
	}
}

func TestHTTPChecker_Check_ServerError(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Status:     "500 Internal Server Error",
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	checker := NewHTTPCheckerWithClient(mockClient)
	endpoint := common.Endpoint{
		URL:     "http://example.com",
		Method:  "GET",
		Timeout: 5 * time.Second,
	}

	result := checker.Check(context.Background(), endpoint)

	if result.Status != common.StatusDown {
		t.Errorf("Expected status Down, got %s", result.Status)
	}

	if result.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code 500, got %d", result.StatusCode)
	}
}

func TestHTTPChecker_Check_NetworkError(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("connection refused")
		},
	}

	checker := NewHTTPCheckerWithClient(mockClient)
	endpoint := common.Endpoint{
		URL:     "http://example.com",
		Method:  "GET",
		Timeout: 5 * time.Second,
	}

	result := checker.Check(context.Background(), endpoint)

	if result.Status != common.StatusUnreachable {
		t.Errorf("Expected status Unreachable, got %s", result.Status)
	}

	if !strings.Contains(result.Error, "connection refused") {
		t.Errorf("Expected error about connection, got: %s", result.Error)
	}
}

func TestHTTPChecker_Check_WithHeaders(t *testing.T) {
	var capturedHeaders http.Header
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			capturedHeaders = req.Header
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	checker := NewHTTPCheckerWithClient(mockClient)
	endpoint := common.Endpoint{
		URL:     "http://example.com",
		Method:  "GET",
		Timeout: 5 * time.Second,
		Headers: map[string]string{
			"Authorization": "Bearer token123",
			"Custom-Header": "value",
		},
	}

	checker.Check(context.Background(), endpoint)

	if capturedHeaders.Get("Authorization") != "Bearer token123" {
		t.Errorf("Authorization header not set correctly")
	}

	if capturedHeaders.Get("Custom-Header") != "value" {
		t.Errorf("Custom-Header not set correctly")
	}
}

func TestHTTPChecker_Check_MustMatchStatus(t *testing.T) {
	tests := []struct {
		name           string
		responseStatus int
		expectedStatus int
		mustMatch      bool
		expectedResult string
	}{
		{"match success", 200, 200, true, common.StatusUp},
		{"match failure", 201, 200, true, common.StatusDegraded},
		{"no strict match", 201, 200, false, common.StatusUp},
		{"500 with strict match", 500, 200, true, common.StatusDown},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: tc.responseStatus,
						Status:     http.StatusText(tc.responseStatus),
						Body:       io.NopCloser(strings.NewReader("")),
					}, nil
				},
			}

			checker := NewHTTPCheckerWithClient(mockClient)
			endpoint := common.Endpoint{
				URL:             "http://example.com",
				Method:          "GET",
				Timeout:         5 * time.Second,
				MustMatchStatus: tc.mustMatch,
				ExpectedStatus:  tc.expectedStatus,
			}

			result := checker.Check(context.Background(), endpoint)

			if result.Status != tc.expectedResult {
				t.Errorf("%s: Expected status %s, got %s", tc.name, tc.expectedResult, result.Status)
			}

			if tc.mustMatch && tc.responseStatus != tc.expectedStatus {
				// Should have UnexpectedStatusCodeMessage
				found := slices.Contains(result.Messages, common.UnexpectedStatusCodeMessage)
				if !found {
					t.Errorf("Expected message about unexpected status code")
				}
			}
		})
	}
}

func TestHTTPChecker_Check_MaxLatency(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Simulate slow response
			time.Sleep(200 * time.Millisecond)
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	checker := NewHTTPCheckerWithClient(mockClient)
	endpoint := common.Endpoint{
		URL:        "http://example.com",
		Method:     "GET",
		Timeout:    5 * time.Second,
		MaxLatency: 100 * time.Millisecond, // Expect degraded due to latency
	}

	result := checker.Check(context.Background(), endpoint)

	if result.Status != common.StatusDegraded {
		t.Errorf("Expected status Degraded due to latency, got %s", result.Status)
	}

	// Check for latency message
	found := slices.Contains(result.Messages, common.UnexpectedLatencyMessage)
	if !found {
		t.Errorf("Expected message about unexpected latency")
	}

	if result.Elapsed < 200 {
		t.Errorf("Expected elapsed time ~200ms, got %dms", result.Elapsed)
	}
}

func TestHTTPChecker_Check_ContextTimeout(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			select {
			case <-req.Context().Done():
				return nil, req.Context().Err()
			case <-time.After(2 * time.Second):
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil
			}
		},
	}

	checker := NewHTTPCheckerWithClient(mockClient)
	endpoint := common.Endpoint{
		URL:     "http://example.com",
		Method:  "GET",
		Timeout: 100 * time.Millisecond, // Short timeout
	}

	result := checker.Check(context.Background(), endpoint)

	if result.Status != common.StatusUnreachable {
		t.Errorf("Expected status Unreachable due to timeout, got %s", result.Status)
	}

	if !strings.Contains(result.Error, "context deadline exceeded") {
		t.Errorf("Expected timeout error, got: %s", result.Error)
	}
}

func TestHTTPChecker_Check_InvalidURL(t *testing.T) {
	checker := NewHTTPChecker() // Using real client for this edge case

	// Create context
	ctx := context.Background()

	// Test with invalid URL
	endpoint := common.Endpoint{
		URL:     "://invalid-url",
		Method:  "GET",
		Timeout: 5 * time.Second,
	}

	result := checker.Check(ctx, endpoint)

	if result.Status != common.StatusUnreachable {
		t.Errorf("Expected status Unreachable for invalid URL, got %s", result.Status)
	}

	if result.Error == "" {
		t.Errorf("Expected error for invalid URL")
	}
}

// Table-driven test for status code classification
func TestHTTPChecker_Check_StatusClassification(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   string
	}{
		{200, common.StatusUp},
		{201, common.StatusUp},
		{299, common.StatusUp},
		{300, common.StatusDegraded},
		{404, common.StatusDegraded},
		{429, common.StatusDegraded},
		{500, common.StatusDown},
		{503, common.StatusDown},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("Status%d", tc.statusCode), func(t *testing.T) {
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: tc.statusCode,
						Status:     http.StatusText(tc.statusCode),
						Body:       io.NopCloser(strings.NewReader("")),
					}, nil
				},
			}

			checker := NewHTTPCheckerWithClient(mockClient)
			endpoint := common.Endpoint{
				URL:     "http://example.com",
				Method:  "GET",
				Timeout: 5 * time.Second,
			}

			result := checker.Check(context.Background(), endpoint)

			if result.Status != tc.expected {
				t.Errorf("Status code %d: Expected %s, got %s",
					tc.statusCode, tc.expected, result.Status)
			}
		})
	}
}
