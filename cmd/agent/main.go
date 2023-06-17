package main

import (
	"metric-alert/internal/agent"
	"metric-alert/internal/server/logger"
)

func main() {
	logger.InitLogger()

	cfg, err := getConfig()
	if err != nil {
		logger.Log.Fatal().Msg("err get config: " + err.Error())
	}

	a := agent.NewAgent(cfg.ReportIntervalSec, cfg.PollIntervalSec, cfg.ServerHost, cfg.SignKey)
	logger.Log.Info().Interface("cfg", cfg).Msg("agent will send reports to " + cfg.ServerHost)

	a.Run()
}
