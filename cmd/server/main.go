package main

import (
	"log"
	"metric-alert/internal/storage"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Fatal("err get config: " + err.Error())
	}

	memStorage := storage.NewMemStore()
	a := NewApp(memStorage)
	log.Default().Println("server start on " + cfg.ServerHost)
	a.Run(cfg.ServerHost)
}
