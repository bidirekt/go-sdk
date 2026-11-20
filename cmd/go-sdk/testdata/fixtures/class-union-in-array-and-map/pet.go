package petstore

type Payment interface{ payment() }

type CardPayment struct {
	Last4 string `json:"last4" bidirekt:"field"`
}

type PixPayment struct {
	Key string `json:"key" bidirekt:"field"`
}

func (CardPayment) payment() {}
func (PixPayment) payment()  {}

type Order struct {
	Payments []Payment          `json:"payments" bidirekt:"field union={card: CardPayment, pix: PixPayment}"`
	ByMethod map[string]Payment `json:"byMethod" bidirekt:"field union={card: CardPayment, pix: PixPayment}"`
	Attempts *[]*Payment        `json:"attempts" bidirekt:"field union={card: CardPayment, pix: PixPayment} optional=true"`
}
