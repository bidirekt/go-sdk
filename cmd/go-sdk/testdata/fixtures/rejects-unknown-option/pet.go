package petstore

type Pet struct {
	ID   int64  `json:"id" bidirekt:"field nmae=petId"`
	Name string `json:"name" bidirekt:"field nullable"`
	Age  int    `json:"age" bidirekt:"field integer"`
}
