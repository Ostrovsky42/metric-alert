package main

import (
	"os"

	"github.com/rs/zerolog"
)

func main() {
	log := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()

	cfg, err := getConfig()
	if err != nil {
		log.Fatal().Msg("err get config: " + err.Error())
	}

	a := NewApp(cfg, log)
	log.Info().Msg("server start on " + cfg.ServerHost)

	a.Run()
}
