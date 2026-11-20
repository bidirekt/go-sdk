package petstore

type Audit struct {
	CreatedBy string `json:"createdBy" bidirekt:"field"`
}

type Timestamps struct {
	Audit
	CreatedAt string `json:"createdAt" bidirekt:"field"`
	UpdatedAt string `json:"updatedAt" bidirekt:"field optional=false"`
}

type Pet struct {
	ID int64 `json:"id" bidirekt:"field"`
	*Timestamps
	Name string `json:"name" bidirekt:"field"`
}

type Owner struct {
	Name string `json:"name" bidirekt:"field"`
	Timestamps
}
