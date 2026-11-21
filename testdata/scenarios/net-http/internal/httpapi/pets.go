package httpapi

import (
	"encoding/json"
	"net/http"

	"example.com/petstore/internal/pets"
	"github.com/bidirekt/go-sdk"
)

type PetHandler struct {
	store *pets.Store
}

func NewPetHandler(store *pets.Store) *PetHandler {
	return &PetHandler{store: store}
}

func (this *PetHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /pets", this.List)
	mux.HandleFunc("POST /pets", this.Create)
	mux.HandleFunc("GET /pets/{id}", this.Get)
	mux.HandleFunc("PATCH /pets/{id}/status", this.ChangeStatus)
	mux.HandleFunc("DELETE /pets/{id}", this.Delete)
}

var _ = bidirekt.Provides("GET", "/pets",
	bidirekt.Response(http.StatusOK, []pets.Pet{}),
)

func (this *PetHandler) List(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, this.store.All())
}

var _ = bidirekt.Provides("POST", "/pets",
	bidirekt.Request(pets.NewPet{}),
	bidirekt.Response(http.StatusCreated, pets.Pet{}),
	bidirekt.Response(http.StatusUnprocessableEntity, ValidationError{}),
)

func (this *PetHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newPet pets.NewPet
	if err := json.NewDecoder(r.Body).Decode(&newPet); err != nil || newPet.Name == "" {
		writeJSON(w, http.StatusUnprocessableEntity, ValidationError{Field: "name", Message: "a pet needs a name"})
		return
	}
	writeJSON(w, http.StatusCreated, this.store.Add(newPet))
}

var _ = bidirekt.Provides("GET", "/pets/*",
	bidirekt.Response(http.StatusOK, pets.Pet{}),
	bidirekt.Response(http.StatusNotFound, nil),
)

func (this *PetHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, valid := pathID(r)
	pet, found := this.store.Find(id)
	if !valid || !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, pet)
}

var _ = bidirekt.Provides("PATCH", "/pets/*/status",
	bidirekt.Request(pets.StatusChange{}),
	bidirekt.Response(http.StatusOK, pets.Pet{}),
	bidirekt.Response(http.StatusNotFound, nil),
	bidirekt.Response(http.StatusUnprocessableEntity, ValidationError{}),
)

func (this *PetHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	var change pets.StatusChange
	if err := json.NewDecoder(r.Body).Decode(&change); err != nil || change.Status == "" {
		writeJSON(w, http.StatusUnprocessableEntity, ValidationError{Field: "status", Message: "a status is required"})
		return
	}
	id, valid := pathID(r)
	pet, found := this.store.ChangeStatus(id, change.Status)
	if !valid || !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, pet)
}

var _ = bidirekt.Provides("DELETE", "/pets/*",
	bidirekt.Response(http.StatusNoContent, nil),
	bidirekt.Response(http.StatusNotFound, nil),
)

func (this *PetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, valid := pathID(r)
	if !valid || !this.store.Remove(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
