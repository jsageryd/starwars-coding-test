package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
)

//go:embed index.tmpl
var index string

var tmpl *template.Template

func init() {
	var err error

	tmpl, err = template.New("index").Parse(index)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}
}

func main() {
	a := API{
		Core: &Core{
			SWAPIBaseURL: "https://swapi.dev/api",
		},
	}

	http.HandleFunc("/", a.ui)
	http.HandleFunc("/top-fat-characters", a.topFatCharacters)

	addr := ":8080"

	log.Printf("Listening at %s...", addr)

	http.ListenAndServe(addr, nil)
}

type Core struct {
	SWAPIBaseURL string
}

func (c *Core) topFatCharacters() ([]Character, error) {
	var characters []Character

	nextURL := strings.TrimRight(c.SWAPIBaseURL, "/") + "/people/"

	for nextURL != "" {
		resp, err := http.Get(nextURL)
		if err != nil {
			return nil, fmt.Errorf("error querying SWAPI: %v", err)
		}

		var respBody struct {
			Next    string      `json:"next"`
			Results []Character `json:"results"`
		}

		if json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			return nil, fmt.Errorf("error reading SWAPI response: %v", err)
		}

		characters = append(characters, respBody.Results...)

		nextURL = respBody.Next
	}

	return topFat(characters, 20), nil
}

type API struct {
	Core *Core
}

func (a *API) ui(w http.ResponseWriter, r *http.Request) {
	characters, err := a.Core.topFatCharacters()
	if err != nil {
		http.Error(w, "unknown error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	data := struct {
		Characters []Character
	}{
		Characters: characters,
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("Error rendering page: %v", err)
		return
	}

	buf.WriteTo(w)
}

func (a *API) topFatCharacters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	characters, err := a.Core.topFatCharacters()
	if err != nil {
		http.Error(w, "unknown error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

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
