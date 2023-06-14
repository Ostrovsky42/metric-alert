package main

import (
	"metric-alert/internal/agent"
	"metric-alert/internal/logger"
)

func main() {
	logger.InitLogger()

	cfg, err := getConfig()
	if err != nil {
		logger.Log.Fatal().Msg("err get config: " + err.Error())
	}

	a := agent.NewAgent(cfg.ReportIntervalSec, cfg.PollIntervalSec, cfg.ServerHost)
	logger.Log.Info().Msg("agent will send reports to " + cfg.ServerHost)

	a.Run()
}
