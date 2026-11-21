package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/bidirekt/go-sdk"
)

type Adoption struct {
	PetID   int64  `json:"petId" bidirekt:"field"`
	Adopter string `json:"adopter" bidirekt:"field"`
}

var adoptions = map[int64]Adoption{1: {PetID: 1, Adopter: "Ada"}}

var _ = bidirekt.Provides("GET", "/adoptions/{petId}",
	bidirekt.Response(http.StatusOK, Adoption{}),
	bidirekt.Response(http.StatusNotFound, nil),
)

func getAdoption(w http.ResponseWriter, r *http.Request) {
	petID, err := strconv.ParseInt(r.PathValue("petId"), 10, 64)
	adoption, found := adoptions[petID]
	if err != nil || !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(adoption)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /adoptions/{petId}", getAdoption)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
