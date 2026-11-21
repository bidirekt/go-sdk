package petstore

type Pet struct {
	ID     int64  `json:"id" bidirekt:"field"`
	Name   string `json:"name" bidirect:"field"`
	Breed  string `json:"breed" Bidirekt:"field"`
	Color  string `json:"color" bidirekt_:"field"`
	Weight int    `json:"weight" BIDIREKT:"field"`
}
