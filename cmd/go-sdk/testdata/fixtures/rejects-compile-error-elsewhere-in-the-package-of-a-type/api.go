package app

import (
	"net/http"

	"example.com/app/models"
	"github.com/bidirekt/go-sdk"
)

type Tag models.Tag

type Ticket struct {
	Tag Tag `json:"tag" bidirekt:"field"`
}

var _ = bidirekt.Provides(http.MethodGet, "/tickets/*",
	bidirekt.Response(http.StatusOK, Ticket{}),
)
