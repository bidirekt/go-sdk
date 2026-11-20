package petstore

type Adopter interface{ adopter() }

type Owner struct {
	Name string `json:"name"`
}

func (Owner) adopter() {}

type Pet struct {
	Name        string            `json:"name" bidirekt:"field"`
	Keeper      string            `json:"keeper" bidirekt:"field type=Person"`
	Friends     []string          `json:"friends" bidirekt:"field type=array items=Friend"`
	Notes       map[string]string `json:"notes" bidirekt:"field type=map values=Note"`
	AdoptedFrom Adopter           `json:"adoptedFrom" bidirekt:"field union={owner: Owner}"`
}
