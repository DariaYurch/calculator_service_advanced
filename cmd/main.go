package main

import (
	"log"
	"os/exec"
	"runtime"
)

func main() {
	var cmd1, cmd2 *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd1 = exec.Command("cmd", "/C", "start", "cmd", "/K", "go run ./orchestrator/main.go")
		cmd2 = exec.Command("cmd", "/C", "start", "cmd", "/K", "go run ./agent/main.go")
	case "linux":
		cmd1 = exec.Command("x-terminal-emulator", "-e", "go run ./service1/main.go")
		cmd2 = exec.Command("x-terminal-emulator", "-e", "go run ./service2/main.go")
	case "darwin":
		cmd1 = exec.Command("osascript", "-e", `tell application "Terminal" to do script "go run ./service1/main.go"`)
		cmd2 = exec.Command("osascript", "-e", `tell application "Terminal" to do script "go run ./service2/main.go"`)
	}

	err := cmd1.Start()
	if err != nil {
		log.Fatalf("Error starting orchestrator/main.go: %v", err)
	}
	err = cmd2.Start()
	if err != nil {
		log.Fatalf("Error starting agent/main.go: %v", err)
	}

}
