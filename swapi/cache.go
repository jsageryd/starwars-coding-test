package swapi

import (
	"sync"

	"github.com/jsageryd/starwars-coding-test/starwars"
)

type cache struct {
	mu         sync.RWMutex
	characters []starwars.Character
}

func newCache() *cache {
	return &cache{}
}

func (c *cache) SetCharacters(cs []starwars.Character) {
	c.mu.Lock()
	c.characters = make([]starwars.Character, len(cs))
	copy(c.characters, cs)
	c.mu.Unlock()
}

func (c *cache) GetCharacters() (cs []starwars.Character, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.characters == nil {
		return nil, false
	}
	cs = make([]starwars.Character, len(c.characters))
	copy(cs, c.characters)
	return cs, true
}
