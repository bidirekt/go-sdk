package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"example.com/storefront/petstore"
)

var page = template.Must(template.New("window").Parse(`<ul>{{range .}}<li>{{.Name}} ({{.Status}})</li>{{end}}</ul>`))

func main() {
	client := &petstore.Client{BaseURL: os.Getenv("PETSTORE_URL"), HTTPClient: &http.Client{Timeout: 5 * time.Second}}
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		pets, err := client.ListPets(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("the pet store is unavailable: %v", err), http.StatusBadGateway)
			return
		}
		_ = page.Execute(w, pets)
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
}
