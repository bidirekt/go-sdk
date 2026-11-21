package main

import (
	"log"
	"net/http"

	"example.com/petstore/internal/httpapi"
	"example.com/petstore/internal/pets"
)

func main() {
	mux := http.NewServeMux()
	httpapi.NewPetHandler(pets.NewStore()).Register(mux)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
