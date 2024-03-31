package main

import (
	"fulli9/platform"
	"fulli9/shared"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	client, database, err := shared.ConnectDB()
	if err != nil {
		log.Fatalf("Error while connecting to mongoDB: %s.\nExiting.", err)
	}
	defer shared.DisConnectDB(client)

	rtr := platform.New(database)

	log.Print("Server listening on http://localhost:3005/")
	if err := http.ListenAndServe("0.0.0.0:3005", rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
