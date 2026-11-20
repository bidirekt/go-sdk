package petstore

type Catalog struct {
	Categories Categories `json:"categories" bidirekt:"field"`
}

type Categories []Category

type Category struct {
	Name     string     `json:"name" bidirekt:"field"`
	Children Categories `json:"children" bidirekt:"field"`
}
