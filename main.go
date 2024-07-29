package main

import (
	"log"
	"net/http"

	"github.com/bignyap/verifyjwt/router"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	mux := http.NewServeMux()

	router.RegisterHandlers(mux)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Error while starting the server")
	}
}
