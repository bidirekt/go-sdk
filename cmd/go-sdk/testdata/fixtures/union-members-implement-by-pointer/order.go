package petstore

type Payment interface{ payment() }

type Card struct {
	Number string `json:"number" bidirekt:"field"`
}

type Pix struct {
	Key string `json:"key" bidirekt:"field"`
}

func (*Card) payment() {}
func (Pix) payment()   {}

type Order struct {
	Payment Payment `json:"payment" bidirekt:"field union={card: Card, pix: Pix}"`
}
