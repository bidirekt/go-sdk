package petstore

type Pet struct {
	Name   string      `json:"name" bidirekt:"field"`
	Extra  any         `json:"extra" bidirekt:"field"`
	Others interface{} `json:"others" bidirekt:"field"`
}
