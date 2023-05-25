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

	fsys, err := fs.Sub(webapp, "dist")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		// Handle API requests here
		fmt.Fprintf(w, "Hello from the API!")
	})

	http.Handle("/", http.FileServer(http.FS(fsys)))

	log.Println("Listening on :3000...")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
