package routes

import (
	"net/http"

	"example.com/app/models"
	"github.com/bidirekt/go-sdk
)

var _ = bidirekt.Provides(http.MethodGet, "/pets/*",
	bidirekt.Response(http.StatusOK, models.Pet{}),
)
