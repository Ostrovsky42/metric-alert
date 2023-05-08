package main

import (
	"flag"
	"metric-alert/internal/agent"
	"net/http"
)

var port = flag.String("a", "8080", "HTTP server endpoint address")
var reportIntervalSec = flag.Int("r", 10, "frequency of sending metrics")
var pollIntervalSec = flag.Int("p", 2, "metric polling frequency")

func main() {
	flag.Parse()

	serverAddress := "http://localhost:" + *port
	a := agent.NewAgent(*reportIntervalSec, *pollIntervalSec, serverAddress)
	a.Run()

	mux := http.NewServeMux()
	err := http.ListenAndServe(`:8081`, mux)
	if err != nil {
		panic(err)
	}
}
