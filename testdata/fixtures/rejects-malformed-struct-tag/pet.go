package petstore

type Pet struct {
	ID    int64  `json:"id" bidirekt:"field"`
	Name  string `json:"b",bidirekt:"field"`
	Breed string `json: "d" bidirekt:"field"`
	Color string `json:"color" xml:"bidirekt"`
}
