package petstore

type Pet struct {
	Weight float64 `json:"weight" bidirekt:"field type=float"`
}
