package petstore

type Pet struct {
	ID   int64  `json:"id" bidirekt:"field"`
	Name string `json:"name" bidirekt:"field"`
}

type NewPet struct {
	Name string `json:"name" bidirekt:"field"`
}

type Pets []Pet

type PetsByName map[string]Pet
