package main

import (
	"embed"
	"flag"
	"os"
	"score/app/runners/confirmer"
	"score/app/runners/preauth"
	"score/app/runners/server"
	"score/app/runners/worker"
)

//go:embed dist/*
var dist embed.FS

type ExecutionMode string

const (
	ServerMode    ExecutionMode = "server"
	WorkerMode    ExecutionMode = "worker"
	ConfirmerMode ExecutionMode = "confirmer"
	PreauthMode   ExecutionMode = "preauth"
)

func main() {
	executionMode := getExecutionMode()
	err := RunApp(executionMode)
	if err != nil {
		panic(err)
	}
}

func RunApp(mode ExecutionMode) error {
	var err error
	switch mode {
	case ServerMode:
		err = server.Run(dist)
	case WorkerMode:
		err = worker.Run()
	case ConfirmerMode:
		err = confirmer.Run()
	case PreauthMode:
		err = preauth.Run()
	default:
		err = server.Run(dist)
	}
	return err
}

// Checks the 'SCORE_MODE' environment variable to determine the execution mode. If it is not set, it will check the 'mode' command line option. If that is not set, it will default to 'server' mode.
func getExecutionMode() ExecutionMode {
	if os.Getenv("SCORE_MODE") != "" {
		return ExecutionMode(os.Getenv("SCORE_MODE"))
	}
	executionModePtr := flag.String("mode", "server", "The execution mode of the application.")
	flag.Parse()
	return ExecutionMode(*executionModePtr)
}
