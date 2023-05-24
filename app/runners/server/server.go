package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

var Run = func(webapp embed.FS) error {
	fmt.Println("Running server...")

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		// Handle API requests here
		fmt.Fprintf(w, "Hello from the API!")
	})

	fsys, err := fs.Sub(webapp, "dist")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(fsys)))

	log.Println("Listening on :8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
