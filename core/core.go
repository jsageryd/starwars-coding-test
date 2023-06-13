package core

import (
	"fmt"
	"sort"
	"strconv"

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

func (c *Core) TopOldCharacters() ([]starwars.Character, error) {
	characters, err := c.swapiClient.People()
	if err != nil {
		return nil, fmt.Errorf("error fetching characters from SWAPI: %v", err)
	}

	return topOld(characters, 20), nil
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

// topOld returns the top N oldest characters according to their birth year.
func topOld(cs []starwars.Character, n int) []starwars.Character {
	var validCs []starwars.Character

	// Skip the ones with a birth year we cannot parse
	for _, c := range cs {
		if _, err := absBirthYear(c.BirthYear); err == nil {
			validCs = append(validCs, c)
		}
	}

	sort.Slice(validCs, func(i, j int) bool {
		a, _ := absBirthYear(validCs[i].BirthYear)
		b, _ := absBirthYear(validCs[j].BirthYear)

		if a == b {
			return validCs[i].Name < validCs[j].Name
		}

		return a < b
	})

	if len(validCs) > n {
		validCs = validCs[:n]
	}

	return validCs
}

func absBirthYear(year string) (float64, error) {
	if len(year) < 4 {
		return 0, fmt.Errorf("unknown birth date format: %s", year)
	}

	numberStr := year[:len(year)-3]
	suffix := year[len(year)-3:]

	number, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return 0, fmt.Errorf("unknown birth date format: %s", year)
	}

	switch suffix {
	case "BBY":
		return -number, nil
	case "ABY":
		return number, nil
	default:
		return 0, fmt.Errorf("unknown birth date format: %s", year)
	}
}
