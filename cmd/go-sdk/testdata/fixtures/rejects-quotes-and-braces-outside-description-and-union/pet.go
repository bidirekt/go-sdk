package petstore

type Pet struct {
	A string `json:"a" bidirekt:"field name='petName'"`
	B string `json:"b" bidirekt:"field type={string}"`
	C string `json:"c" bidirekt:"field description=\"Shown\""`
	D string `json:"d" bidirekt:"field name=\"petName\""`
}
