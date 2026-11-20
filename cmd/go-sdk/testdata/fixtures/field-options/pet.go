package petstore

type Adopter interface{ adopter() }

type Owner struct {
	Name string `json:"name" bidirekt:"field"`
}

type Shelter struct {
	City string `json:"city" bidirekt:"field"`
}

func (Owner) adopter()   {}
func (Shelter) adopter() {}

type Pet struct {
	ID          int64             `json:"id" bidirekt:"field name=petId"`
	Weight      float64           `json:"weight" bidirekt:"field type=float"`
	Keeper      string            `json:"keeper" bidirekt:"field type=Owner"`
	Tags        []string          `json:"tags" bidirekt:"field type=array items=Owner"`
	Notes       map[string]string `json:"notes" bidirekt:"field type=map values=integer"`
	Nickname    string            `json:"nickname" bidirekt:"field optional=true"`
	Breed       *string           `json:"breed" bidirekt:"field nullable=false"`
	Color       string            `json:"color" bidirekt:"field description=Shown"`
	AdoptedFrom Adopter           `json:"adoptedFrom" bidirekt:"field union={person: Owner, shelter: Shelter}"`
}
