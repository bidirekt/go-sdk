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
	Payments []Payment          `json:"payments" bidirekt:"field union={card: Card, pix: Pix}"`
	ByMethod map[string]Payment `json:"byMethod" bidirekt:"field union={card: Card, pix: Pix}"`
	Backup   *Payment           `json:"backup" bidirekt:"field union={card: Card, pix: Pix}"`
}
