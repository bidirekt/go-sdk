package main

import (
	"log"
	"net/http"

	"example.com/bookshelf/authors"
	"example.com/bookshelf/books"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	bookHandler := books.NewHandler()
	authorHandler := authors.NewHandler()

	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.Logger, middleware.Recoverer)

	router.Route(books.BasePath, func(router chi.Router) {
		router.Get("/", bookHandler.List)
		router.Post("/", bookHandler.Create)
		router.Route("/{isbn}", func(router chi.Router) {
			router.Get("/", bookHandler.Get)
			router.Put("/", bookHandler.Replace)
			router.Delete("/", bookHandler.Delete)
		})
	})

	router.Group(func(router chi.Router) {
		router.Use(middleware.BasicAuth("bookshelf", map[string]string{"admin": "secret"}))
		router.Get(authors.BasePath, authorHandler.List)
		router.Post(authors.BasePath, authorHandler.Create)
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}
