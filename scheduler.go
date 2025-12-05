package main

import (
	"context"
	"fmt"
	"time"
)

type Scheduler struct {
	endpoints []Endpoint
	checkers  map[string]Checker // "http" → HTTPChecker, etc.
	results   chan Result
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

func (s *Scheduler) runEndpoint(ep Endpoint) {
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
				// Send an error result to maintain consistency
				s.results <- Result{
					URL:       ep.URL,
					Status:    StatusUnreachable,
					Timestamp: time.Now(),
					Error:     "no checker registered for type",
					Message:   fmt.Sprintf("Checker for type %q not found", ep.Type),
				}
				continue
			}
			res := checker.Check(context.Background(), ep)
			s.results <- res
		case <-s.stop: // in this case we stop
			return
		}
	}
}
func (s *Scheduler) Stop() {
	close(s.stop) // Signal all goroutines to stop
}
