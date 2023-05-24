package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

//go:embed dist/*
var dist embed.FS

var ExecutionMode string = ""

func main() {
	parseCommandLineOptions()
	fmt.Println("Execution Mode: ", ExecutionMode)
}

func parseCommandLineOptions() {
	modeOptionPointer := flag.String("mode", "server", "the mode to run the program in")
	flag.Parse()
	if modeOptionPointer == nil {
		ExecutionMode = "server"
	} else {
		ExecutionMode = *modeOptionPointer
	}

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		// Handle API requests here
		fmt.Fprintf(w, "Hello from the API!")
	})

	fsys, err := fs.Sub(dist, "dist")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(fsys)))

	log.Println("Listening on :8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
