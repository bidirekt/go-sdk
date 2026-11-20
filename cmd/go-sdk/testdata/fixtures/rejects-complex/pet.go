package petstore

type Pet struct {
	Name  string     `json:"name" bidirekt:"field"`
	Phase complex128 `json:"phase" bidirekt:"field"`
}
