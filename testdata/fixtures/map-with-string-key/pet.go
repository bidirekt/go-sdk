package petstore

type Breed string

type Labels map[string]int

type Pet struct {
	Notes   map[string]string `json:"notes" bidirekt:"field"`
	ByBreed map[Breed]*Pet    `json:"byBreed" bidirekt:"field"`
	Labels  Labels            `json:"labels" bidirekt:"field"`
}
