package main

import (
	"log"
	"metric-alert/internal/server/config"
	"metric-alert/internal/server/logger"
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof
)

func main() {
	logger.InitLogger()
	cfg := config.GetConfig()

	a := NewApp(cfg)
	defer a.Close()
	logger.Log.Info().Interface("cfg", cfg).Msg("server start on " + cfg.ServerHost)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	a.Run()
}
