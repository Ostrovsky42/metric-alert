package main

import (
	"metric-alert/internal/logger"
)

func main() {
	logger.InitLogger()

	cfg := getConfig()
	a := NewApp(cfg)
	defer a.Close()
	logger.Log.Info().Msg("server start on " + cfg.ServerHost)

	a.Run()
}
