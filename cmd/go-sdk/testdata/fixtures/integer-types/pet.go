package petstore

type Pet struct {
	A int    `json:"a" bidirekt:"field"`
	B int8   `json:"b" bidirekt:"field"`
	C int16  `json:"c" bidirekt:"field"`
	D int32  `json:"d" bidirekt:"field"`
	E int64  `json:"e" bidirekt:"field"`
	F uint   `json:"f" bidirekt:"field"`
	G uint8  `json:"g" bidirekt:"field"`
	H uint16 `json:"h" bidirekt:"field"`
	I uint32 `json:"i" bidirekt:"field"`
	J uint64 `json:"j" bidirekt:"field"`
}
