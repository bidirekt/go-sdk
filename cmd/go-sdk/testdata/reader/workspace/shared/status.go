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

var _ = bidirekt.Provides("GET", "/shared",
	bidirekt.Response(200, nil),
)

type InvoiceStatus string

const (
	InvoiceStatusOpen InvoiceStatus = "open"
	InvoiceStatusPaid InvoiceStatus = "paid"
)

var _ = bidirekt.Enum(InvoiceStatusOpen, InvoiceStatusPaid)
