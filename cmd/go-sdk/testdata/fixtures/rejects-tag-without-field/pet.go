package petstore

type Pet struct {
	ID    int64  `json:"id" bidirekt:"name=petId"`
	Name  string `json:"name" bidirekt:""`
	Breed string `json:"breed" bidirekt:"fields"`
}
