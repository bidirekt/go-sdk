package orders

import (
	"example.com/shop/shared"
	"github.com/bidirekt/go-sdk"
)

type Status string

const (
	StatusPlaced    Status = "placed"
	StatusPaid      Status = "paid"
	StatusShipped   Status = "shipped"
	StatusCancelled Status = "cancelled"
)

var _ = bidirekt.Enum(StatusPlaced, StatusPaid, StatusShipped, StatusCancelled)

type Payment interface {
	payment()
}

type CardPayment struct {
	Brand    string `json:"brand" bidirekt:"field"`
	LastFour string `json:"lastFour" bidirekt:"field"`
}

type BankTransfer struct {
	IBAN      string `json:"iban" bidirekt:"field"`
	Reference string `json:"reference" bidirekt:"field"`
}

func (CardPayment) payment() {}

func (BankTransfer) payment() {}

type Line struct {
	SKU      string `json:"sku" bidirekt:"field"`
	Quantity int    `json:"quantity" bidirekt:"field"`
	Cents    int64  `json:"cents" bidirekt:"field description='Unit price in cents'"`
}

type Order struct {
	ID         string  `json:"id" bidirekt:"field"`
	CustomerID string  `json:"customerId" bidirekt:"field"`
	Status     Status  `json:"status" bidirekt:"field"`
	Lines      []Line  `json:"lines" bidirekt:"field"`
	Payment    Payment `json:"payment" bidirekt:"field union={card: CardPayment, bank: BankTransfer}"`
	shared.Timestamps
}

type PlaceOrder struct {
	CustomerID string  `json:"customerId" bidirekt:"field"`
	Lines      []Line  `json:"lines" bidirekt:"field"`
	Payment    Payment `json:"payment" bidirekt:"field union={card: CardPayment, bank: BankTransfer}"`
}
