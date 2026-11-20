package petstore

type Timestamps struct {
	CreatedAt string `json:"createdAt" bidirekt:"field"`
	UpdatedAt string `json:"updatedAt" bidirekt:"field"`
}

type Pet struct {
	ID int64 `json:"id" bidirekt:"field"`
	*Timestamps
	Name string `json:"name" bidirekt:"field"`
}
