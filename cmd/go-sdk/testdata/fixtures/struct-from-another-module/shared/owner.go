package shared

type Owner struct {
	Name  string `json:"name" bidirekt:"field"`
	Email string `json:"email"`
}

type Unused struct {
	Name string `json:"name" bidirekt:"field"`
}
