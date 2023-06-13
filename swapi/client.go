package swapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jsageryd/starwars-coding-test/starwars"
)

type Client struct {
	baseURL string
	cache   *cache
}

func NewClient(baseURL string) *Client {
	return &Client{
		cache:   newCache(),
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (c *Client) People() ([]starwars.Character, error) {
	if cs, ok := c.cache.GetCharacters(); ok {
		return cs, nil
	}

	var characters []starwars.Character

	nextURL := c.baseURL + "/people/"

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

	c.cache.SetCharacters(characters)

	return characters, nil
}
