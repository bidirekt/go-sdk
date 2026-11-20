package petstore

type PetCard struct {
	Owner        Owner   `json:"owner" bidirekt:"field"`
	FormerOwners []Owner `json:"formerOwners" bidirekt:"field"`
	Address      Address `json:"address" bidirekt:"field"`
}
