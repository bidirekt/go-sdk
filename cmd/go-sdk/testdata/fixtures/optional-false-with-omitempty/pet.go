package petstore

type Pet struct {
	Name string `json:"name,omitempty" bidirekt:"field optional=false"`
	Age  int    `json:"age" bidirekt:"field optional=true"`
}
