package main

import (
	"log"
	"net/http"

	"github.com/jsageryd/starwars-coding-test/api"
	"github.com/jsageryd/starwars-coding-test/core"
	"github.com/jsageryd/starwars-coding-test/swapi"
)

func main() {
	mux := http.NewServeMux()

	swapiClient := swapi.NewClient("https://swapi.dev/api")

	api.New(
		core.New(
			swapiClient,
		),
	).Register(mux)

	go func() {
		log.Printf("Warming up cache...")
		if _, err := swapiClient.People(); err != nil {
			log.Printf("Error warming up cache: %v", err)
		} else {
			log.Printf("Cache warm")
		}
	}()

	addr := ":8080"

	log.Printf("Listening at %s...", addr)

	http.ListenAndServe(addr, mux)
}
