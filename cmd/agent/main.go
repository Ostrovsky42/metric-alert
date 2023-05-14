package main

import (
	"log"
	"metric-alert/internal/agent"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Fatal("err get config: " + err.Error())
	}

	a := agent.NewAgent(cfg.ReportIntervalSec, cfg.PollIntervalSec, cfg.ServerHost)
	log.Default().Println("agent will send reports to " + cfg.ServerHost)

	a.Run()
}
