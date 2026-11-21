package main

import (
	"log"
	"net/http"

	"example.com/acme/apps/api/catalogapi"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", catalogapi.NewServer().Routes()))
}
