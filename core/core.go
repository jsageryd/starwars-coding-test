package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/jsageryd/starwars-coding-test/starwars"
)

type Core struct {
	swapiBaseURL string
}

func New(swapiBaseURL string) *Core {
	return &Core{
		swapiBaseURL: strings.TrimRight(swapiBaseURL, "/"),
	}
}

func (c *Core) TopFatCharacters() ([]starwars.Character, error) {
	var characters []starwars.Character

	nextURL := c.swapiBaseURL + "/people/"

	for nextURL != "" {
		resp, err := http.Get(nextURL)
		if err != nil {
			return nil, fmt.Errorf("error querying SWAPI: %v", err)
		}

		var respBody struct {
			Next    string               `json:"next"`
			Results []starwars.Character `json:"results"`
		}

		if json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			return nil, fmt.Errorf("error reading SWAPI response: %v", err)
		}

		characters = append(characters, respBody.Results...)

		nextURL = respBody.Next
	}

	return topFat(characters, 20), nil
}

// topFat returns the top N fattest characters according to their BMI.
func topFat(cs []starwars.Character, n int) []starwars.Character {
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
