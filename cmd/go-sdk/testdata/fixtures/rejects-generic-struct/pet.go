package petstore

type Page[T any] struct {
	Items []T `json:"items" bidirekt:"field"`
	Total int `json:"total" bidirekt:"field"`
}

type Pet struct {
	Name string `json:"name" bidirekt:"field"`
}

type Catalog struct {
	Pets Page[Pet] `json:"pets" bidirekt:"field"`
}
