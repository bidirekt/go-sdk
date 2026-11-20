package petstore

type Pet struct {
	ID     int64    `json:"id,string" bidirekt:"field"`
	Alive  bool     `json:"alive,string" bidirekt:"field"`
	Weight float64  `json:"weight,string" bidirekt:"field"`
	Rank   *int     `json:"rank,string" bidirekt:"field"`
	Tags   []string `json:"tags,string" bidirekt:"field"`
}
