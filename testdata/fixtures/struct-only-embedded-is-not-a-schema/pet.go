package petstore

type Timestamps struct {
	CreatedAt string `json:"createdAt" bidirekt:"field"`
	UpdatedAt string `json:"updatedAt" bidirekt:"field"`
}

type Pet struct {
	Name string `json:"name" bidirekt:"field"`
}

type PetResponse struct {
	Pet
	Timestamps
}
