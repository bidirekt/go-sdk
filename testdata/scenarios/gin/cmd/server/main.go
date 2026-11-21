package main

import (
	"log"

	"example.com/ticketing/internal/api"
	"example.com/ticketing/internal/tickets"
)

func main() {
	router := api.NewRouter(tickets.NewService())
	log.Fatal(router.Run(":8080"))
}
