package main

import (
	"log"
	"metric-alert/internal/agent"
)

func main() {
	parseFlags()

	a := agent.NewAgent(reportIntervalSec, pollIntervalSec, host)
	log.Default().Println("agent will send reports to " + host)
	a.Run()
}
