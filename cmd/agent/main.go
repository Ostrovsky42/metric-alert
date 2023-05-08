package main

import (
	"metric-alert/internal/agent"
	"net/http"
	"time"
)

func main() {
	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second
	serverURL := "http://localhost:8080"

	a := agent.NewAgent(reportInterval, pollInterval, serverURL)
	a.Run()

	mux := http.NewServeMux()
	err := http.ListenAndServe(`:8081`, mux)
	if err != nil {
		panic(err)
	}
}
