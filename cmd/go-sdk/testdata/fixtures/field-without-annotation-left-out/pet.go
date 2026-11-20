package petstore

type Pet struct {
	ID            float64 `json:"id" bidirekt:"field"`
	Name          string  `json:"name" bidirekt:"field"`
	InternalNotes string  `json:"internalNotes"`
}
