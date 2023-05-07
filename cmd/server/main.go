package main

import "metric-alert/internal/storage"

func main() {
	memStorage := storage.NewMemStore()
	a := NewApp(memStorage)
	a.Run()
}
