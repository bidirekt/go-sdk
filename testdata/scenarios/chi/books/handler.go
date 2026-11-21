package books

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/bidirekt/go-sdk"
	"github.com/go-chi/chi/v5"
)

const BasePath = "/v1/books"

type Handler struct {
	mutex sync.RWMutex
	books map[string]Book
}

func NewHandler() *Handler {
	return &Handler{books: map[string]Book{}}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeProblem(w http.ResponseWriter, status int, detail string) {
	writeJSON(w, status, Problem{Title: http.StatusText(status), Status: status, Detail: detail})
}

var _ = bidirekt.Provides(http.MethodGet, BasePath,
	bidirekt.Response(http.StatusOK, []Book{}),
)

func (this *Handler) List(w http.ResponseWriter, r *http.Request) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	all := make([]Book, 0, len(this.books))
	for _, book := range this.books {
		all = append(all, book)
	}
	writeJSON(w, http.StatusOK, all)
}

var _ = bidirekt.Provides(http.MethodPost, BasePath,
	bidirekt.Request(BookInput{}),
	bidirekt.Response(http.StatusCreated, Book{}),
	bidirekt.Response(http.StatusBadRequest, Problem{}),
	bidirekt.Response(http.StatusConflict, Problem{}),
)

func (this *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input BookInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeProblem(w, http.StatusBadRequest, err.Error())
		return
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if _, exists := this.books[input.ISBN]; exists {
		writeProblem(w, http.StatusConflict, "a book with this ISBN exists")
		return
	}
	book := bookFromInput(input)
	this.books[book.ISBN] = book
	writeJSON(w, http.StatusCreated, book)
}

var _ = bidirekt.Provides(http.MethodGet, BasePath+"/*",
	bidirekt.Response(http.StatusOK, Book{}),
	bidirekt.Response(http.StatusNotFound, Problem{}),
)

func (this *Handler) Get(w http.ResponseWriter, r *http.Request) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	book, found := this.books[chi.URLParam(r, "isbn")]
	if !found {
		writeProblem(w, http.StatusNotFound, "")
		return
	}
	writeJSON(w, http.StatusOK, book)
}

var _ = bidirekt.Provides(http.MethodPut, BasePath+"/*",
	bidirekt.Request(BookInput{}),
	bidirekt.Response(http.StatusOK, Book{}),
	bidirekt.Response(http.StatusBadRequest, Problem{}),
	bidirekt.Response(http.StatusNotFound, Problem{}),
)

func (this *Handler) Replace(w http.ResponseWriter, r *http.Request) {
	var input BookInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeProblem(w, http.StatusBadRequest, err.Error())
		return
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	isbn := chi.URLParam(r, "isbn")
	if _, found := this.books[isbn]; !found {
		writeProblem(w, http.StatusNotFound, "")
		return
	}
	input.ISBN = isbn
	book := bookFromInput(input)
	this.books[isbn] = book
	writeJSON(w, http.StatusOK, book)
}

var _ = bidirekt.Provides(http.MethodDelete, BasePath+"/*",
	bidirekt.Response(http.StatusNoContent, nil),
	bidirekt.Response(http.StatusNotFound, Problem{}),
)

func (this *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	isbn := chi.URLParam(r, "isbn")
	if _, found := this.books[isbn]; !found {
		writeProblem(w, http.StatusNotFound, "")
		return
	}
	delete(this.books, isbn)
	w.WriteHeader(http.StatusNoContent)
}

func bookFromInput(input BookInput) Book {
	authors := make([]Author, 0, len(input.AuthorIDs))
	for _, id := range input.AuthorIDs {
		authors = append(authors, Author{ID: id})
	}
	return Book{ISBN: input.ISBN, Title: input.Title, Authors: authors, Price: input.Price, Pages: input.Pages}
}
