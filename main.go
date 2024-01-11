package main

import (
	"fulli9/platform"
	"log"
	"net/http"
)

func main() {
	rtr := platform.New()

	log.Print("Server listening on http://localhost:3000/")
	if err := http.ListenAndServe("0.0.0.0:3000", rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
