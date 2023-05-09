package main

import "flag"

var host string

func parseFlags() {
	flag.StringVar(&host, "a", "localhost:8080", "server endpoint host")
	flag.Parse()
}
