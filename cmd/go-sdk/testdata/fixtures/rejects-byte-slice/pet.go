package petstore

type Image []byte

type Pet struct {
	Name      string `json:"name" bidirekt:"field"`
	Photo     []byte `json:"photo" bidirekt:"field"`
	Thumbnail Image  `json:"thumbnail" bidirekt:"field"`
}
