package petstore

type Pet struct {
	Name     string         `json:"name" bidirekt:"field"`
	ByNumber map[int]string `json:"byNumber" bidirekt:"field"`
}
