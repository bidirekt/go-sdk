package petstore

type Pet struct {
	Tags []*string `json:"tags" bidirekt:"field"`
}
