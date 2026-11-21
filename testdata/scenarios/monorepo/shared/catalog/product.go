package catalog

import (
	"example.com/acme/shared/money"
	"github.com/bidirekt/go-sdk"
)

type Category string

const (
	CategoryBooks       Category = "books"
	CategoryElectronics Category = "electronics"
	CategoryGarden      Category = "garden"
)

var _ = bidirekt.Enum(CategoryBooks, CategoryElectronics, CategoryGarden)

type Availability string

const (
	AvailabilityInStock    Availability = "in_stock"
	AvailabilityBackorder  Availability = "backorder"
	AvailabilityDiscounted Availability = "discontinued"
)

var _ = bidirekt.Enum(AvailabilityInStock, AvailabilityBackorder, AvailabilityDiscounted)

type Product struct {
	SKU        string            `json:"sku" bidirekt:"field"`
	Name       string            `json:"name" bidirekt:"field"`
	Category   Category          `json:"category" bidirekt:"field"`
	Price      money.Amount      `json:"price" bidirekt:"field"`
	Attributes map[string]string `json:"attributes,omitempty" bidirekt:"field"`
}

type StockLevel struct {
	SKU       string `json:"sku" bidirekt:"field"`
	Available int    `json:"available" bidirekt:"field"`
	Reserved  int    `json:"reserved" bidirekt:"field"`
}
