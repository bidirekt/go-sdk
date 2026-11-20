package petstore

type Owner struct {
	Name string `json:"name" bidirekt:"field"`
}

type Pet struct {
	Owner  Owner  `json:"owner,omitzero" bidirekt:"field"`
	Scores [3]int `json:"scores,omitzero" bidirekt:"field"`
	Name   string `json:"name,omitzero" bidirekt:"field"`
}
