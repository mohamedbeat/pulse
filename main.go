package main

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
	if bufferSize < 10 {
		bufferSize = 10
	}

	scheduler := Scheduler{
		endpoints: config.Endpoints,
		checkers: map[string]Checker{
			HTTPType: httpChecker,
		},
		results: make(chan Result, bufferSize),
		stop:    make(chan struct{}),
	}

	//Starting scheduler
	go scheduler.Start()

	//Getting scheduler results
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
