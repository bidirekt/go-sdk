package money

type Money struct {
	Cents int `json:"cents" bidirekt:"field"`
}

func (Money) MarshalJSON() ([]byte, error) {
	return []byte(`"0.00"`), nil
}
