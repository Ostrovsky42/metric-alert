package main

import (
	"metric-alert/internal/server/config"
	"metric-alert/internal/server/logger"
)

func main() {
	logger.InitLogger()

	cfg := config.GetConfig()

	a := NewApp(cfg)
	defer a.Close()
	logger.Log.Info().Interface("cfg", cfg).Msg("server start on " + cfg.ServerHost)

	a.Run()
}
