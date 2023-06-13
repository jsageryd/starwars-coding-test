package starwars

type Character struct {
	Name      string  `json:"name"`
	Height    float64 `json:"height,string,omitempty"` // height in cm
	Mass      float64 `json:"mass,string,omitempty"`   // mass in kg
	BirthYear string  `json:"birth_year,omitempty"`    // birth year, "<year> BBY" or "<year> ABY"
}
