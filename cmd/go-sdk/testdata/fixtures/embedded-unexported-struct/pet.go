package petstore

type timestamps struct {
	CreatedAt string `json:"createdAt" bidirekt:"field"`
}

type Pet struct {
	timestamps
	Name string `json:"name" bidirekt:"field"`
}
