package shared

import (
	"github.com/bidirekt/go-sdk"
)

type PetStatus string

const (
	PetStatusAvailable PetStatus = "available"
	PetStatusAdopted   PetStatus = "adopted"
)

var _ = bidirekt.Enum(PetStatusAvailable, PetStatusAdopted)

type InvoiceStatus string

const (
	InvoiceStatusOpen InvoiceStatus = "open"
	InvoiceStatusPaid InvoiceStatus = "paid"
)

var _ = bidirekt.Enum(InvoiceStatusOpen, InvoiceStatusPaid)

type Pet struct {
	Name   string    `json:"name" bidirekt:"field"`
	Status PetStatus `json:"status" bidirekt:"field"`
}
