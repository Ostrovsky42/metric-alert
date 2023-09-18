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
	printBuildInfo()

	cfg := config.GetConfig()

	a := NewApp(cfg)
	defer a.Close()
	logger.Log.Info().Interface("cfg", cfg).Msg("server start on " + cfg.ServerHost)

	go func() {
		log.Println(http.ListenAndServe(cfg.ProfilerHost, nil))
	}()

	a.Run()
}
