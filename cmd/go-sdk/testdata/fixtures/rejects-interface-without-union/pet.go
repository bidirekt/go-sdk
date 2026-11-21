package petstore

type Adopter interface{ adopter() }

type Pet struct {
	Name        string  `json:"name" bidirekt:"field"`
	AdoptedFrom Adopter `json:"adoptedFrom" bidirekt:"field"`
}
