package main

import (
	"os"

	"github.com/rs/zerolog"
	"metric-alert/internal/agent"
)

func main() {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg, err := getConfig()
	if err != nil {
		log.Fatal().Msg("err get config: " + err.Error())
	}

	a := agent.NewAgent(cfg.ReportIntervalSec, cfg.PollIntervalSec, cfg.ServerHost, log)
	log.Info().Msg("agent will send reports to " + cfg.ServerHost)

	a.Run()
}
