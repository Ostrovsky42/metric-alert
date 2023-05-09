package main

import (
	"metric-alert/internal/storage"
)

func main() {
	parseFlags()

	memStorage := storage.NewMemStore()
	a := NewApp(memStorage)
	a.Run(host)
}
