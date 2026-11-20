package petstore

type Pet struct {
	Name     string `json:"name" bidirekt:"field"`
	Nickname string `json:"nickname" bidirekt:"field name=name"`
	Breed    string `json:"breed" bidirekt:"field"`
}
