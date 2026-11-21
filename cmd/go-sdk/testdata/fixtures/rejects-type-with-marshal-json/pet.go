package petstore

type Color struct {
	Red int `json:"red" bidirekt:"field"`
}

func (Color) MarshalJSON() ([]byte, error) {
	return []byte(`"red"`), nil
}

type Pet struct {
	Name  string `json:"name" bidirekt:"field"`
	Color Color  `json:"color" bidirekt:"field"`
}
