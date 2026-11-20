package petstore

type Pet struct {
	Nickname *string `json:"nickname" bidirekt:"field optional=true"`
}
