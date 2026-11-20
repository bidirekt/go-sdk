package petstore

type Pet struct {
	ID        int64  `json:"id" bidirekt:"field"`
	name      string `json:"name" bidirekt:"field"`
	CreatedBy string `json:"-" bidirekt:"field"`
}
