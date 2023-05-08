package main

import (
	"flag"
	"metric-alert/internal/storage"
)

var serverAddress = flag.String("a", "localhost:8080", "HTTP server endpoint address")

func main() {
	flag.Parse()

	memStorage := storage.NewMemStore()
	a := NewApp(memStorage)
	a.Run(*serverAddress)
}
