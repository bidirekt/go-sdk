package petstore

type Pet struct {
	Weight float64 `json:"weight" bidirekt:"field description='In kilograms'"`
	Age    int     `json:"age" bidirekt:"field description='Age in years since 2001' optional=true"`
}
