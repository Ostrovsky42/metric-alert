package main

import "flag"

var host string
var reportIntervalSec int
var pollIntervalSec int

func parseFlags() {
	flag.StringVar(&host, "a", "localhost:8080", "server endpoint address")
	flag.IntVar(&reportIntervalSec, "r", 10, "frequency of sending metrics")
	flag.IntVar(&pollIntervalSec, "p", 2, "metric polling frequency")

	flag.Parse()
}
