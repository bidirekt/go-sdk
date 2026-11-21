package petstore

type List[T any] []T

type Index[T any] map[string]T

type Owner struct {
	Name string `json:"name" bidirekt:"field"`
}

type Pet struct {
	Name    string       `json:"name" bidirekt:"field"`
	Friends List[Pet]    `json:"friends" bidirekt:"field"`
	Owners  Index[Owner] `json:"owners" bidirekt:"field"`
}
