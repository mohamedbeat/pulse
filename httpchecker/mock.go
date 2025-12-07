package httpchecker

import (
	"io"
	"net/http"
	"strings"
	"time"
)

// MockHTTPClient is a mock HTTP client for testing
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       io.NopCloser(strings.NewReader("")),
	}, nil
}

// Test helpers to create mock responses
func createMockResponse(statusCode int, body string, latency time.Duration) *http.Response {
	// Simulate latency if needed
	if latency > 0 {
		time.Sleep(latency)
	}

	return &http.Response{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}
