package petstore

type Pet struct {
	Weight float64 `json:"weight" bidirekt:"field"`
	Height float32 `json:"height" bidirekt:"field"`
}
