package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"metric-alert/internal/agent"
	"metric-alert/internal/agent/config"
	"metric-alert/internal/server/logger"
)

func main() {
	logger.InitLogger()

	cfg := config.GetConfig()

	a := agent.NewAgent(cfg)
	logger.Log.Info().Interface("cfg", cfg).Msg("agent will send reports to " + cfg.ServerHost)

	go func() {
		log.Println(http.ListenAndServe(cfg.ServerHost, nil))
	}()

	a.Run()
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
