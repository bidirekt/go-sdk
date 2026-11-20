package petstore

type Pet struct {
	Scores map[string]float64 `json:"scores" bidirekt:"field"`
}
