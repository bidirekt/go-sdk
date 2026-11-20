package petstore

type Animal struct {
	Name string  `json:"name" bidirekt:"field"`
	Age  float64 `json:"age" bidirekt:"field"`
}

type Dog struct {
	Name string `json:"name" bidirekt:"field description='The name the owner calls the dog'"`
	Animal
}
