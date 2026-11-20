package petstore

type Pet struct {
	Nickname string `json:"nickname" bidirekt:"field description='Shown on the card'"`

	// The name the owner calls the pet
	Name string `json:"name" bidirekt:"field"`
}
