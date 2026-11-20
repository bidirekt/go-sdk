package petstore

type Pet struct {
	Name   string      `json:"name" bidirekt:"field"`
	Events chan string `json:"events" bidirekt:"field"`
}
