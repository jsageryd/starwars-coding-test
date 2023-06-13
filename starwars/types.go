package starwars

type Character struct {
	Name   string  `json:"name"`
	Height float64 `json:"height,string"` // height in cm
	Mass   float64 `json:"mass,string"`   // mass in kg
}
