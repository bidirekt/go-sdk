package petstore

type Pet struct {
	Tags    []string  `json:"tags" bidirekt:"field"`
	Scores  []float64 `json:"scores" bidirekt:"field"`
	Aliases [2]string `json:"aliases" bidirekt:"field"`
}
