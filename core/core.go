package core

import (
	"fmt"
	"sort"

	"github.com/jsageryd/starwars-coding-test/starwars"
	"github.com/jsageryd/starwars-coding-test/swapi"
)

type Core struct {
	swapiClient *swapi.Client
}

func New(client *swapi.Client) *Core {
	return &Core{
		swapiClient: client,
	}
}

func (c *Core) TopFatCharacters() ([]starwars.Character, error) {
	characters, err := c.swapiClient.People()
	if err != nil {
		return nil, fmt.Errorf("error fetching characters from SWAPI: %v", err)
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
