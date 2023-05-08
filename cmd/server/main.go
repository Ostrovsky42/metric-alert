package main

import (
	"flag"
	"metric-alert/internal/storage"
)

var port = flag.String("a", "8080", "HTTP server endpoint address")

func main() {
	flag.Parse()

	memStorage := storage.NewMemStore()
	a := NewApp(memStorage)
	a.Run(*port)
}
