package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
)

func main() {
	a := API{
		SWAPIBaseURL: "https://swapi.dev/api",
	}

	http.HandleFunc("/top-fat-characters", a.topFatCharacters)

	addr := ":8080"

	log.Printf("Listening at %s...", addr)

	http.ListenAndServe(addr, nil)
}

type API struct {
	SWAPIBaseURL string
}

func (a *API) topFatCharacters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var characters []Character

	nextURL := strings.TrimRight(a.SWAPIBaseURL, "/") + "/people/"

	for nextURL != "" {
		resp, err := http.Get(nextURL)
		if err != nil {
			http.Error(w, "unknown error", http.StatusInternalServerError)
			log.Printf("Error querying SWAPI: %v", err)
			return
		}

		var respBody struct {
			Next    string      `json:"next"`
			Results []Character `json:"results"`
		}

		if json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			http.Error(w, "unknown error", http.StatusInternalServerError)
			log.Printf("Error reading SWAPI response: %v", err)
			return
		}

		characters = append(characters, respBody.Results...)

		nextURL = respBody.Next
	}

	characters = topFat(characters, 20)

	w.Header().Set("content-type", "application/json")

	json.NewEncoder(w).Encode(&characters)
}

type Character struct {
	Name   string  `json:"name"`
	Height float64 `json:"height,string"` // height in cm
	Mass   float64 `json:"mass,string"`   // mass in kg
}

// topFat returns the top N fattest characters according to their BMI.
func topFat(cs []Character, n int) []Character {
	sort.Slice(cs, func(i, j int) bool {
		bmi := func(height, mass float64) float64 {
			heightCm := height / 100
			return mass / (heightCm * heightCm)
		}

		return bmi(cs[i].Height, cs[i].Mass) > bmi(cs[j].Height, cs[j].Mass)
	})

	if len(cs) > n {
		cs = cs[:n]
	}

	return cs
}
