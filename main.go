package main

import (
	"net/http"
	"time"
)

func main() {
	ep := Endpoint{
		Name:            "first",
		URL:             "http://localhost:9000/health",
		Method:          MethodGet,
		Timeout:         1 * time.Second,
		Interval:        2 * time.Second,
		ExpectedStatus:  http.StatusOK,
		MustMatchStatus: true,
		MaxLatency:      20 * time.Millisecond,
		Type:            "http",
		// Headers: ,
		// BodyContains: ,
		// BodyRegex: ,

	}
	ep2 := Endpoint{
		Name:            "second",
		URL:             "http://localhost:9000/latency",
		Method:          MethodGet,
		Timeout:         1 * time.Second,
		Interval:        10 * time.Second,
		ExpectedStatus:  http.StatusOK,
		MustMatchStatus: true,
		MaxLatency:      20 * time.Millisecond,
		Type:            "http",
		// Headers: ,
		// BodyContains: ,
		// BodyRegex: ,

	}
	httpChecker := NewHTTPChecker()
	eps := []Endpoint{ep, ep2}
	scheduler := Scheduler{
		endpoints: eps,
		checkers: map[string]Checker{
			"http": httpChecker,
		},
		results: make(chan Result),
		stop:    make(chan struct{}),
	}
	go scheduler.Start()

	for result := range scheduler.results {

		switch result.Status {
		case StatusDown, StatusUnreachable:
			Error("",
				"url", result.URL,
				"status", result.Status,
				"status_code", result.StatusCode,
				"timestamp", result.Timestamp,
				"error", result.Error,
				"message", result.Message,
				"elapsed", result.Elapsed,
			)
		case StatusDegraded:
			Warn("saving_result",
				"url", result.URL,
				"status", result.Status,
				"status_code", result.StatusCode,
				"timestamp", result.Timestamp,
				"error", result.Error,
				"message", result.Message,
				"elapsed", result.Elapsed,
			)
		default:
			Info("saving_result",
				"url", result.URL,
				"status", result.Status,
				"status_code", result.StatusCode,
				"timestamp", result.Timestamp,
				"error", result.Error,
				"message", result.Message,
				"elapsed", result.Elapsed,
			)
		}

		// store.Save(result)

		// oldStatus := alertState[result.EndpointID]
		// newStatus := result.Status

		// if oldStatus != newStatus {
		// alertState[result.EndpointID] = newStatus
		// alert := Alert{...}
		// for _, n := range m.notifiers {
		//     go n.Notify(context.Background(), alert) // async!
		// }
		// }
	}
}
