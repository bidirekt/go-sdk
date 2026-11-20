package petstore

type Owner struct {
	ID int64 `json:"id" bidirekt:"field"`
}

type Shelter struct {
	ID int64 `json:"id" bidirekt:"field"`
}

type Adoption struct {
	Owner
	Shelter
	Date string `json:"date" bidirekt:"field"`
}
