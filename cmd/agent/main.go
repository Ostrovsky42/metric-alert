package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"metric-alert/internal/agent"
	"metric-alert/internal/agent/config"
	"metric-alert/internal/server/logger"
)

func main() {
	logger.InitLogger()
	printBuildInfo()

	cfg := config.GetConfig()

	a := agent.NewAgent(cfg)
	logger.Log.Info().Interface("cfg", cfg).Msg("agent will send reports to " + cfg.ServerHost)

	go func() {
		log.Println(http.ListenAndServe(cfg.ProfilerHost, nil))
	}()

	a.Run()
}
