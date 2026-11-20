package petstore

type Pet struct {
	PetName string `bidirekt:"field"`
	Breed   string `json:",omitempty" bidirekt:"field"`
	Color   string `json:"color" bidirekt:"field name=colour"`
	Owner   string `json:"owner,omitempty" bidirekt:"field name=keeper"`
	Odd     string `json:"a\\b" bidirekt:"field"`
}
