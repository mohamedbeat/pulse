package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	Info("Initializing ...")

	configPath := ParseFlags()
	Info("configFile path", "path", configPath)

	config, err := LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	Debug("Globals", "Globals", config.Globals)
	Debug("Config", "config", config)

	httpChecker := NewHTTPChecker()

	// Buffer size: at least 10, or 2x the number of endpoints (whichever is larger)
	// This handles bursts when multiple endpoints complete checks simultaneously
	bufferSize := len(config.Endpoints) * 2
	bufferSize = max(bufferSize, 10)

	scheduler := Scheduler{
		endpoints: config.Endpoints,
		checkers: map[string]Checker{
			HTTPType: httpChecker,
		},
		results: make(chan Result, bufferSize),
		stop:    make(chan struct{}),
	}

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// A goroutine that will close the results channel when scheduler stops
	// This will break the for-loop below
	go func() {
		<-quit // Wait for shutdown signal
		Info("Shutdown signal received")

		// Stop the scheduler (this closes s.stop channel)
		scheduler.Stop()

		// Wait for all endpoint goroutines to finish
		time.Sleep(1 * time.Second)

		// Close the results channel to break the for-loop
		close(scheduler.results)
	}()

	//Starting scheduler
	go scheduler.Start()

	//Getting scheduler results
	for result := range scheduler.results {
		switch result.Status {
		case StatusDown, StatusUnreachable:
			Error("Error",
				"url", result.URL,
				"status", result.Status,
				"status_code", result.StatusCode,
				"timestamp", result.Timestamp,
				"error", result.Error,
				"message", result.Message,
				"elapsed", result.Elapsed,
			)
		case StatusDegraded:
			Warn("Warning",
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

		// Info("Shutdown complete")

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
