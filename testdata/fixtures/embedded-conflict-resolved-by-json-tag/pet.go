package petstore

type Code struct {
	Code string `json:"Code" bidirekt:"field"`
}

type Label struct {
	Code string `bidirekt:"field"`
	Text string `json:"text" bidirekt:"field"`
}

type Pet struct {
	Code
	Label
}
