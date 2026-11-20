package pets

type Adopter interface{ adopter() }

type Owner struct {
	Name string `json:"name" bidirekt:"field"`
}

type Shelter struct {
	City string `json:"city" bidirekt:"field"`
}

func (Owner) adopter() {}

type Pet struct {
	Name        string  `json:"name" bidirekt:"field"`
	AdoptedFrom Adopter `json:"adoptedFrom" bidirekt:"field union={person: Owner, shelter: Shelter}"`
}
