package petstore

type Names []string

type Pet struct {
	Tags      []string `json:"tags" bidirekt:"field"`
	Scores    [3]int   `json:"scores" bidirekt:"field"`
	Nicknames Names    `json:"nicknames" bidirekt:"field"`
	Grid      [][]int  `json:"grid" bidirekt:"field"`
}
