package main

import (
	"calculator_advanced_v2/agent/worker"
	"calculator_advanced_v2/orchestrator/services/storage"
	"os"
	"strconv"
)

func main() {
	storage.InitConnectionDB()
	computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil || computingPower <= 0 {
		computingPower = 1
	}
	for i := 0; i < computingPower; i++ {
		go worker.StartWorker(i)
	}
	select {}
}
