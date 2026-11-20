package petstore

type Owner struct {
	Name string `json:"name" bidirekt:"field"`
}

type Pet struct {
	A string            `json:"a" bidirekt:"field optional=yes"`
	B string            `json:"b" bidirekt:"field nullable=1"`
	C string            `json:"c" bidirekt:"field name="`
	D []Owner           `json:"d" bidirekt:"field type=[]Owner"`
	E map[string]Owner  `json:"e" bidirekt:"field type=map[string]Owner"`
	F []Owner           `json:"f" bidirekt:"field type=array items=array"`
	G map[string]string `json:"g" bidirekt:"field type=map values=map<string>"`
	H string            `json:"h" bidirekt:"field type=array"`
	I string            `json:"i" bidirekt:"field type=map"`
	J string            `json:"j" bidirekt:"field items=string"`
	K string            `json:"k" bidirekt:"field values=string"`
	L any               `json:"l" bidirekt:"field type=Owner union={a: Owner}"`
	M any               `json:"m" bidirekt:"field union=Owner"`
	N any               `json:"n" bidirekt:"field union={a: Owner, a: Owner}"`
	O any               `json:"o" bidirekt:"field union={a Owner}"`
	P string            `json:"p" bidirekt:"field description='Not closed"`
}
