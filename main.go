package main

import (
	"embed"
	"flag"
	"score/app"
)

//go:embed dist/*
var dist embed.FS

var ExecutionMode string = ""

func main() {
	parseCommandLineOptions()
	switch ExecutionMode {
	case "server":
		app.RunServer(dist)
	case "worker":
		app.RunWorker()
	case "confirmer":
		app.RunPostAccountConfirmationHandler()
	case "preauth":
		app.RunPreAuthenticationHandler()
	default:
		app.RunServer(dist)
	}

}

func parseCommandLineOptions() {
	modeOptionPointer := flag.String("mode", "server", "the mode to run the program in")
	flag.Parse()
	if modeOptionPointer == nil {
		ExecutionMode = "server"
	} else {
		ExecutionMode = *modeOptionPointer
	}
}
