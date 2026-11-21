package books

type Author struct {
	ID   int    `json:"id" bidirekt:"field"`
	Name string `json:"name" bidirekt:"field"`
}

type Book struct {
	ISBN      string   `json:"isbn" bidirekt:"field"`
	Title     string   `json:"title" bidirekt:"field"`
	Authors   []Author `json:"authors" bidirekt:"field"`
	Price     float64  `json:"price" bidirekt:"field type=float description='In euros'"`
	Pages     int      `json:"pages,omitempty" bidirekt:"field"`
	Published bool     `json:"published" bidirekt:"field"`
	Series    *string  `json:"series" bidirekt:"field"`
}

type BookInput struct {
	ISBN      string  `json:"isbn" bidirekt:"field"`
	Title     string  `json:"title" bidirekt:"field"`
	AuthorIDs []int   `json:"authorIds" bidirekt:"field"`
	Price     float64 `json:"price" bidirekt:"field type=float"`
	Pages     int     `json:"pages,omitempty" bidirekt:"field"`
}

type Problem struct {
	Title  string `json:"title" bidirekt:"field"`
	Status int    `json:"status" bidirekt:"field"`
	Detail string `json:"detail,omitempty" bidirekt:"field"`
}
