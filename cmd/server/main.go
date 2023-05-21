package main

import (
	"os"

	"github.com/rs/zerolog"
	"metric-alert/internal/storage"
)

func main() {
	log := zerolog.New(os.Stdout).With().Caller().Logger()

	cfg, err := getConfig()
	if err != nil {
		log.Fatal().Msg("err get config: " + err.Error())
	}

	memStorage := storage.NewMemStore()
	a := NewApp(memStorage, cfg.ServerHost, log)
	log.Info().Msg("server start on " + cfg.ServerHost)

	a.Run()
}
