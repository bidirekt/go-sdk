package petstore

type Timestamps struct {
	CreatedAt string `json:"createdAt" bidirekt:"field"`
	UpdatedAt string `json:"updatedAt" bidirekt:"field"`
}

type Pet struct {
	Timestamps
	CreatedAt int64 `json:"createdAt" bidirekt:"field"`
}
