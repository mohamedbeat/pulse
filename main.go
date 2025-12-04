package main

import ()

func main() {

	Info("Initializing ...")
	config, err := LoadConfig("")
	if err != nil {
		panic(err)
	}

	Debug("Globals", "Globals", config.Globals)
	Debug("Config", "config", config)

	httpChecker := NewHTTPChecker()

	scheduler := Scheduler{
		endpoints: config.Endpoints,
		checkers: map[string]Checker{
			HTTPType: httpChecker,
		},
		results: make(chan Result),
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
