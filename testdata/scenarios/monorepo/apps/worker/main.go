package main

import (
	"log"
	"time"

	"example.com/acme/shared/catalog"
)

func main() {
	levels := map[string]catalog.StockLevel{}
	for range time.Tick(time.Minute) {
		for sku, level := range levels {
			if level.Available-level.Reserved <= 0 {
				log.Printf("%s is out of stock", sku)
			}
		}
	}
}
