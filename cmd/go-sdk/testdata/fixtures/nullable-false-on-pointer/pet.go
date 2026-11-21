package petstore

type Pet struct {
	Breed *string `json:"breed" bidirekt:"field nullable=false"`
	Name  string  `json:"name" bidirekt:"field nullable=true"`
}
