package petstore

type Pet struct {
	Breed *string `json:"breed,omitempty" bidirekt:"field"`
	Age   *int    `json:"age,omitzero" bidirekt:"field"`
}
