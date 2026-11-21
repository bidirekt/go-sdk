package petstore

type Pet struct {
	Nickname string `json:"nickname,omitempty" bidirekt:"field"`
}
