package petstore

type Key string

type PetsByName map[string]Pet

type Pet struct {
	Name string `json:"name" bidirekt:"field"`
}

type Shelter struct {
	Capacity map[string]float64 `json:"capacity" bidirekt:"field"`
	Notes    map[Key]*string    `json:"notes" bidirekt:"field"`
	Pets     PetsByName         `json:"pets" bidirekt:"field"`
	Tags     map[Key]string     `json:"tags" bidirekt:"field"`
}
