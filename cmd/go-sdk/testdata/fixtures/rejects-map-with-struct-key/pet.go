package petstore

type Owner struct {
	Name string `json:"name" bidirekt:"field"`
}

type Pet struct {
	Name    string           `json:"name" bidirekt:"field"`
	ByOwner map[Owner]string `json:"byOwner" bidirekt:"field"`
}
