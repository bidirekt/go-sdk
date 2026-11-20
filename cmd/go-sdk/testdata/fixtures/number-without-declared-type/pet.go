package petstore

type Pet struct {
	Weight float64 `json:"weight" bidirekt:"field"`
	Age    float64 `json:"age" bidirekt:"field type=integer"`
	Height float64 `json:"height" bidirekt:"field type=float"`
}
