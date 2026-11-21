package petstore

type Pet struct {
	Nickname *string `json:"nickname" bidirekt:"field nullable=false"`
	Name     string  `json:"name,omitempty" bidirekt:"field optional=false"`
}
