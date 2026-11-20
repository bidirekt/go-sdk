package petstore

type Pet struct {
	ID    int64  `json:"id" bidirekt:"field"`
	Name  string `json:"b",bidirekt:"field"`
	Breed string `json: "d" bidirekt:"field"`
	Color string `json:"color" xml:"bidirekt"`
	Owner string `json:"owner" bidirekt:"field
description='Who'"`
	Notes string `json:"notes" bidirekt:"field description='a\qb'"`
}
