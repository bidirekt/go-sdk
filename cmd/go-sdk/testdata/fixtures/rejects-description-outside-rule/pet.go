package petstore

type Pet struct {
	Weight float64 `json:"weight" bidirekt:"field description='In kilograms.'"`
	Age    int     `json:"age" bidirekt:"field description=''"`
	Breed  string  `json:"breed" bidirekt:"field description="`
	Color  string  `json:"color" bidirekt:"field description='Two  spaces'"`
	Name   string  `json:"name" bidirekt:"field description=' leading'"`
	Owner  string  `json:"owner" bidirekt:"field description='Dono do cão'"`
	Diet   string  `json:"diet" bidirekt:"field description='Fish & rice <wet>'"`
}
