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

	api.New(
		core.New(
			swapi.NewClient("https://swapi.dev/api"),
		),
	).Register(mux)

	addr := ":8080"

	log.Printf("Listening at %s...", addr)

	http.ListenAndServe(addr, mux)
}
