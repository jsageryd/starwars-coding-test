package swapi

import "github.com/jsageryd/starwars-coding-test/starwars"

type peopleResponse struct {
	Count   int                  `json:"count"`
	Results []starwars.Character `json:"results"`
}
