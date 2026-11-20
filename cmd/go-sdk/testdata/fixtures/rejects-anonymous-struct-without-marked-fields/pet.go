package petstore

type Pet struct {
	Name    string `json:"name" bidirekt:"field"`
	Address struct {
		City string `json:"city"`
	} `json:"address" bidirekt:"field"`
}
