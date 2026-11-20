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
	Payment Payment  `json:"payment" bidirekt:"field union={card: CardPayment, pix: PixPayment}"`
	Refund  *Payment `json:"refund" bidirekt:"field union={pix: PixPayment, card: CardPayment} optional=true"`
}
