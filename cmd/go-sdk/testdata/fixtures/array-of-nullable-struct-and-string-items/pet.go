package petstore

type Pet struct {
	Friends []*Pet    `json:"friends" bidirekt:"field"`
	Tags    []*string `json:"tags" bidirekt:"field"`
}
