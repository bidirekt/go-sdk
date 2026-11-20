package petstore

type Size int

func (*Size) MarshalText() ([]byte, error) {
	return []byte("small"), nil
}

type Pet struct {
	Name string `json:"name" bidirekt:"field"`
	Size Size   `json:"size" bidirekt:"field"`
}
