package petstore

type Pet struct {
	Notes map[string]*string `json:"notes" bidirekt:"field"`
}
