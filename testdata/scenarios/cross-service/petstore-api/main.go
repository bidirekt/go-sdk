package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"example.com/petstore-api/pets"
	"github.com/bidirekt/go-sdk"
)

type server struct {
	mutex sync.Mutex
	pets  []pets.Pet
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

var _ = bidirekt.Provides(http.MethodGet, "/pets",
	bidirekt.Response(http.StatusOK, []pets.Pet{}),
)

func (this *server) listPets(w http.ResponseWriter, r *http.Request) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	writeJSON(w, http.StatusOK, this.pets)
}

var _ = bidirekt.Provides(http.MethodPost, "/pets",
	bidirekt.Request(pets.NewPet{}),
	bidirekt.Response(http.StatusCreated, pets.Pet{}),
	bidirekt.Response(http.StatusUnprocessableEntity, nil),
)

func (this *server) createPet(w http.ResponseWriter, r *http.Request) {
	var newPet pets.NewPet
	if err := json.NewDecoder(r.Body).Decode(&newPet); err != nil || newPet.Name == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	pet := pets.Pet{ID: int64(len(this.pets) + 1), Name: newPet.Name, Status: pets.StatusAvailable, BirthDate: newPet.BirthDate}
	this.pets = append(this.pets, pet)
	writeJSON(w, http.StatusCreated, pet)
}

var _ = bidirekt.Provides(http.MethodGet, "/pets/*",
	bidirekt.Response(http.StatusOK, pets.Pet{}),
	bidirekt.Response(http.StatusNotFound, nil),
)

func (this *server) getPet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if err != nil || id < 1 || id > int64(len(this.pets)) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, this.pets[id-1])
}

func main() {
	api := &server{}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /pets", api.listPets)
	mux.HandleFunc("POST /pets", api.createPet)
	mux.HandleFunc("GET /pets/{id}", api.getPet)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
