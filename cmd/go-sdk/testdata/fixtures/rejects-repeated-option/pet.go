package petstore

type Pet struct {
	ID int64 `json:"id" bidirekt:"field name=petId name=id"`
}
