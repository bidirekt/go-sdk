package authors

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/bidirekt/go-sdk"
)

const BasePath = "/v1/authors"

type Author struct {
	ID        int      `json:"id" bidirekt:"field"`
	Name      string   `json:"name" bidirekt:"field"`
	Biography string   `json:"biography,omitempty" bidirekt:"field"`
	Languages []string `json:"languages" bidirekt:"field"`
}

type NewAuthor struct {
	Name      string   `json:"name" bidirekt:"field"`
	Biography string   `json:"biography,omitempty" bidirekt:"field"`
	Languages []string `json:"languages" bidirekt:"field"`
}

var _ = bidirekt.Schema(Author{}, bidirekt.Named("AuthorProfile"))

type Handler struct {
	mutex   sync.Mutex
	authors []Author
}

func NewHandler() *Handler {
	return &Handler{}
}

var _ = bidirekt.Provides(http.MethodGet, BasePath,
	bidirekt.Response(http.StatusOK, []Author{}),
	bidirekt.Response(http.StatusUnauthorized, nil),
)

func (this *Handler) List(w http.ResponseWriter, r *http.Request) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(this.authors)
}

var _ = bidirekt.Provides(http.MethodPost, BasePath,
	bidirekt.Request(NewAuthor{}),
	bidirekt.Response(http.StatusCreated, Author{}),
	bidirekt.Response(http.StatusBadRequest, nil),
	bidirekt.Response(http.StatusUnauthorized, nil),
)

func (this *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input NewAuthor
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	author := Author{ID: len(this.authors) + 1, Name: input.Name, Biography: input.Biography, Languages: input.Languages}
	this.authors = append(this.authors, author)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(author)
}
