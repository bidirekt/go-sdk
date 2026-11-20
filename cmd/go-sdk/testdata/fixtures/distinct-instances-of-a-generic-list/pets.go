package petstore

type List[T any] []T

type Index[T any] map[string]T

type Owner struct {
	Name string `json:"name" bidirekt:"field"`
}

type Pet struct {
	Name    string               `json:"name" bidirekt:"field"`
	Friends List[Pet]            `json:"friends" bidirekt:"field"`
	Owners  List[Owner]          `json:"owners" bidirekt:"field"`
	Litters List[List[Pet]]      `json:"litters" bidirekt:"field"`
	Scores  List[List[int]]      `json:"scores" bidirekt:"field"`
	Notes   Index[Index[string]] `json:"notes" bidirekt:"field"`
}
