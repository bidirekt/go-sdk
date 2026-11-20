package petstore

type Owner struct {
	Name string `json:"name" bidirekt:"field"`
}

type Pet struct {
	Name   string            `json:"name,omitempty" bidirekt:"field"`
	Alive  bool              `json:"alive,omitempty" bidirekt:"field"`
	Age    int               `json:"age,omitempty" bidirekt:"field"`
	Weight float64           `json:"weight,omitempty" bidirekt:"field"`
	Tags   []string          `json:"tags,omitempty" bidirekt:"field"`
	Notes  map[string]string `json:"notes,omitempty" bidirekt:"field"`
	Scores [3]int            `json:"scores,omitempty" bidirekt:"field"`
	Owner  Owner             `json:"owner,omitempty" bidirekt:"field"`
}
