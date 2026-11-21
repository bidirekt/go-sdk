package petstore

type Owner struct {
	Name string `json:"name" bidirekt:"field"`
}

type Address struct {
	Street string `json:"street" bidirekt:"field"`
}
