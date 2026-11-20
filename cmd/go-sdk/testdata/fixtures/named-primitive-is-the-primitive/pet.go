package petstore

type PetID int64

type Status string

type Score float64

type Flag bool

type Pet struct {
	ID     PetID  `json:"id" bidirekt:"field"`
	Status Status `json:"status" bidirekt:"field"`
	Score  Score  `json:"score" bidirekt:"field"`
	Alive  Flag   `json:"alive" bidirekt:"field"`
}
