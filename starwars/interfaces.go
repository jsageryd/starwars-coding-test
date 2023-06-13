package starwars

type Core interface {
	TopFatCharacters() ([]Character, error)
}
