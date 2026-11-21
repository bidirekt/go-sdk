package petstore

type Money struct {
	Cents int `json:"cents" bidirekt:"field"`
}

func (Money) MarshalJSON() ([]byte, error) {
	return []byte(`"0.00"`), nil
}

type Code struct {
	Value string `json:"value" bidirekt:"field"`
}

func (*Code) MarshalText() ([]byte, error) {
	return []byte("code"), nil
}
