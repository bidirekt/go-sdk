package petstore

type Pet struct {
	Name     string  `json:"name" bidirekt:"field"`
	Mother   *Pet    `json:"mother" bidirekt:"field"`
	Children []Pet   `json:"children" bidirekt:"field"`
	Owner    *Person `json:"owner" bidirekt:"field"`
}

type Person struct {
	Pets []Pet `json:"pets" bidirekt:"field"`
}
