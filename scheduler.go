package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamedbeat/pulse/common"
)

type Scheduler struct {
	endpoints []common.Endpoint
	checkers  map[string]Checker // "HTTP" → HTTPChecker, etc.
	results   chan common.Result
	stop      chan struct{}
}

func (s *Scheduler) Start() {
	// Validate all endpoints have checkers before starting
	for _, ep := range s.endpoints {
		if _, ok := s.checkers[ep.Type]; !ok {
			Error("missing_checker",
				"endpoint", ep.Name,
				"url", ep.URL,
				"type", ep.Type,
				"message", "No checker registered for endpoint type",
			)
		}
	}

	for _, ep := range s.endpoints {
		go s.runEndpoint(ep) // one goroutine per endpoint
	}
}

func (s *Scheduler) runEndpoint(ep common.Endpoint) {
	ticker := time.NewTicker(ep.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checker, ok := s.checkers[ep.Type]
			if !ok {
				Error("missing_checker",
					"endpoint", ep.Name,
					"url", ep.URL,
					"type", ep.Type,
					"message", "No checker registered for endpoint type, skipping check",
				)

				messages := make([]string, 0)
				messages = append(messages, fmt.Sprintf("Checker for type %q not found", ep.Type))

				// Send an error result to maintain consistency
				s.results <- common.Result{
					URL:       ep.URL,
					Status:    common.StatusUnreachable,
					Timestamp: time.Now(),
					Error:     "no checker registered for type",
					Messages:  messages,
				}
				continue
			}

			res := checker.Check(context.Background(), ep)

			// check results for retry
			if res.Status != common.StatusUp && ep.RetryCounter > 0 {
				Warn("Bad result", res)
				Warn("Retrying")
				ep.LastResult = &res
				ep.RetryCounter -= 1
				continue
			}

			// Resetting RetryCounter and LastResult
			ep.RetryCounter = ep.Retry
			ep.LastResult = nil

			s.results <- res
		case <-s.stop: // in this case we stop
			return
		}
	}
}
func (s *Scheduler) Stop() {
	close(s.stop) // Signal all goroutines to stop
}
