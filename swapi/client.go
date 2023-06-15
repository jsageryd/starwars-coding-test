package swapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

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

	pr, err := c.fetchPeople(1)
	if err != nil {
		return nil, err
	}

	if len(pr.Results) == 0 {
		return nil, nil
	}

	characters := pr.Results

	pageCount := pr.Count/len(pr.Results) + 1

	type ret struct {
		pr  peopleResponse
		err error
	}

	in := make(chan int, pageCount-1)
	out := make(chan ret)

	for page := 2; page <= pageCount; page++ {
		in <- page
	}
	close(in)

	workers := 10
	if workers > pageCount-1 {
		workers = pageCount
	}

	var wg sync.WaitGroup
	wg.Add(workers)

	for n := 0; n < workers; n++ {
		go func() {
			defer wg.Done()

			for page := range in {
				pr, err := c.fetchPeople(page)
				out <- ret{pr, err}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for r := range out {
		if r.err != nil {
			return nil, r.err
		}

		characters = append(characters, r.pr.Results...)
	}

	c.cache.SetCharacters(characters)

	return characters, nil
}

func (c *Client) fetchPeople(page int) (peopleResponse, error) {
	pageURL := fmt.Sprintf("%s/people/?page=%d", c.baseURL, page)

	resp, err := http.Get(pageURL)
	if err != nil {
		return peopleResponse{}, fmt.Errorf("error querying SWAPI: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return peopleResponse{}, fmt.Errorf("SWAPI returned HTTP %d", resp.StatusCode)
	}

	var pr peopleResponse

	if json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return peopleResponse{}, fmt.Errorf("error reading SWAPI response: %v", err)
	}

	return pr, nil
}
