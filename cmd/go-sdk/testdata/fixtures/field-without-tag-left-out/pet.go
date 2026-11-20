package petstore

type Pet struct {
	ID       int64  `json:"id" bidirekt:"field"`
	Name     string `json:"name" bidirekt:"field"`
	Internal string `json:"internal"`
	Secret   string
}
