package main

import (
	"flag"
	"metric-alert/internal/agent"
	"net/http"
)

var port = flag.String("a", "8080", "HTTP server endpoint address")
var reportInterval = flag.Int("r", 10, "input image file")
var pollInterval = flag.Int("p", 2, "input image file")

func main() {
	flag.Parse()

	serverAddress := "http://localhost:" + *port
	a := agent.NewAgent(*reportInterval, *pollInterval, serverAddress)
	a.Run()

	mux := http.NewServeMux()
	err := http.ListenAndServe(`:8081`, mux)
	if err != nil {
		panic(err)
	}
}
