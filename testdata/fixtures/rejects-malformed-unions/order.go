package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type Payment interface{ payment() }

type Card struct {
	Number string `json:"number" bidirekt:"field"`
}

type Pix struct {
	Key string `json:"key" bidirekt:"field"`
}

func (*Card) payment() {}
func (Pix) payment()   {}

type PetStatus string

const (
	PetStatusAvailable PetStatus = "available"
	PetStatusAdopted   PetStatus = "adopted"
)

var _ = bidirekt.Enum(PetStatusAvailable)

type Boleto struct {
	Code string `json:"code" bidirekt:"field"`
}

type Order struct {
	A Payment `json:"a" bidirekt:"field union={card: Card, boleto: Boleto}"`
	B Payment `json:"b" bidirekt:"field union={card: Card, status: PetStatus}"`
	C Payment `json:"c" bidirekt:"field union={card: Card, credit: Card}"`
	D string  `json:"d" bidirekt:"field union={card: Card, pix: Pix}"`
}
