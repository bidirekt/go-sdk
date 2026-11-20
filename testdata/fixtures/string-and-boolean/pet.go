package petstore

type Pet struct {
	Name       string `json:"name" bidirekt:"field"`
	Vaccinated bool   `json:"vaccinated" bidirekt:"field"`
}
