package petstore

type Pet struct {
	Name    string `json:"name" bidirekt:"field"`
	OnAdopt func() `json:"onAdopt" bidirekt:"field"`
}
