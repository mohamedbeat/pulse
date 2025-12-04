package main

import (
	"context"
	"time"
)

type Scheduler struct {
	endpoints []Endpoint
	checkers  map[string]Checker // "http" → HTTPChecker, etc.
	results   chan Result
	stop      chan struct{}
}

func (s *Scheduler) Start() {

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
			checker := s.checkers[ep.Type]
			res := checker.Check(context.Background(), ep)
			s.results <- res
		case <-s.stop:
			return
		}
	}
}
