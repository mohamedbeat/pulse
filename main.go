package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mohamedbeat/pulse/common"
	"github.com/mohamedbeat/pulse/httpchecker"
	"github.com/mohamedbeat/pulse/store"
	"github.com/mohamedbeat/pulse/store/repository"
)

func main() {
	ctx := context.Background()
	Info("Initializing ...")
	envs, err := LoadEnv()
	if err != nil {
		panic(err)
	}

	err = store.Connect("./db.sqlite")
	if err != nil {
		panic(err)
	}
	db := store.GetDBConnection()
	if db == nil {
		panic("Nil db connection")
	}

	// Initializing repository
	HTTPCheckRepo := repository.NewHTTPCheckRepository(db)
	fmt.Println(envs)
	r, err := HTTPCheckRepo.GetLatest(ctx, "latency")
	if err != nil {
		panic(err)
	}
	fmt.Println(r)

	configPath := ParseFlags()
	Info("configFile path", "path", configPath)

	config, err := LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	Debug("Globals", "Globals", config.Globals)
	Debug("Config", "config", config)

	httpChecker := httpchecker.NewHTTPChecker()

	// Buffer size: at least 10, or 2x the number of endpoints (whichever is larger)
	// This handles bursts when multiple endpoints complete checks simultaneously
	bufferSize := len(config.Endpoints) * 2
	bufferSize = max(bufferSize, 10)

	scheduler := Scheduler{
		endpoints: config.Endpoints,
		checkers: map[string]Checker{
			common.HTTPType: httpChecker,
		},
		results: make(chan common.Result, bufferSize),
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
		fmt.Println("messages", result.Messages)
		HTTPCheck := store.HTTPCheck{
			Name:            result.Endpoint.Name,
			Method:          result.Endpoint.Method,
			URL:             result.Endpoint.URL,
			Interval:        result.Endpoint.Interval,
			Timeout:         result.Endpoint.Timeout,
			ExpectedStatus:  result.Endpoint.ExpectedStatus,
			Headers:         result.Endpoint.Headers,
			MustMatchStatus: result.Endpoint.MustMatchStatus,

			//result related
			Status:       result.Status,
			StatusCode:   result.StatusCode,
			ErrorMessage: result.Error,
			Duration:     time.Duration(result.Elapsed),
			CheckedAt:    result.Timestamp,
		}
		fmt.Println("duration: ", HTTPCheck.Duration.Microseconds())
		err := HTTPCheckRepo.Save(ctx, &HTTPCheck)
		if err != nil {
			fmt.Println(err)

		}

		switch result.Status {
		case common.StatusDown, common.StatusUnreachable:
			Error("Error",
				"url", result.URL,
				"status", result.Status,
				"status_code", result.StatusCode,
				"timestamp", result.Timestamp,
				"error", result.Error,
				"messages", result.Messages,
				"elapsed", result.Elapsed,
			)
		case common.StatusDegraded:
			Warn("Warning",
				"url", result.URL,
				"status", result.Status,
				"status_code", result.StatusCode,
				"timestamp", result.Timestamp,
				"error", result.Error,
				"messages", result.Messages,
				"elapsed", result.Elapsed,
			)
		default:
			Info("saving_result",
				"url", result.URL,
				"status", result.Status,
				"status_code", result.StatusCode,
				"timestamp", result.Timestamp,
				"error", result.Error,
				"messages", result.Messages,
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
