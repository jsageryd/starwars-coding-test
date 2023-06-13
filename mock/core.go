package mock

import "github.com/jsageryd/starwars-coding-test/starwars"

type Core struct {
	TopFatCharactersFunc func() ([]starwars.Character, error)
	TopOldCharactersFunc func() ([]starwars.Character, error)
}

func (c *Core) TopFatCharacters() ([]starwars.Character, error) {
	return c.TopFatCharactersFunc()
}

func (c *Core) TopOldCharacters() ([]starwars.Character, error) {
	return c.TopOldCharactersFunc()
}
