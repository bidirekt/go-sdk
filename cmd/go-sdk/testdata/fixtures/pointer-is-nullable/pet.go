package petstore

type Owner struct {
	Name string `json:"name" bidirekt:"field"`
}

type Pet struct {
	Breed *string `json:"breed" bidirekt:"field"`
	Age   *int    `json:"age" bidirekt:"field"`
	Owner *Owner  `json:"owner" bidirekt:"field"`
}
