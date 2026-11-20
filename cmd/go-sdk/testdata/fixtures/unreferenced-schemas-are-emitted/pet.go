package petstore

type Pet struct {
	Name string `json:"name" bidirekt:"field"`
}

type PetCard struct {
	Name string `json:"name" bidirekt:"field"`
}
