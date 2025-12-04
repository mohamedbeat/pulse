package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	ep := Endpoint{
		Name:            "first",
		URL:             "http://localhost:9000/health",
		Method:          MethodGet,
		Timeout:         1 * time.Second,
		Interval:        5 * time.Second,
		ExpectedStatus:  http.StatusOK,
		MustMatchStatus: true,
		MaxLatency:      20 * time.Millisecond,
		Type:            "http",
		// Headers: ,
		// BodyContains: ,
		// BodyRegex: ,

	}
	httpChecker := NewHTTPChecker()
	eps := []Endpoint{ep}
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
		log.Println("saving result", result.URL, result.StatusCode, result.StatusCode, result.Timestamp, result.Error, result.Message, result.ResponseTime)
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

// func main() {
// 	fmt.Println("Pulse !!!!!")
//
// 	ctx := context.Background()
// 	ep := Endpoint{
// 		URL:             "http://localhost:9000/slow",
// 		Method:          MethodGet,
// 		Timeout:         1 * time.Second,
// 		ExpectedStatus:  http.StatusOK,
// 		MustMatchStatus: false,
// 	}
// 	client := NewHTTPChecker()
// 	result, err := client.Check(ctx, ep)
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Println(result)
//
// }
