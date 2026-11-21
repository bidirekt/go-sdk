package shared

const Version = "1"

type Owner struct {
	Name string `json:"name" bidirekt:"field"`
}

type Friend struct {
	Name string `json:"name" bidirekt:"field"`
}
