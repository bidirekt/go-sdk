package app

import (
	"net/http"

	"example.com/app/codes"
	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Provides(http.MethodPost, "/tickets",
	bidirekt.Response(codes.Created, nil),
	bidirekt.Response(codes.Accepted, nil),
)

var _ = bidirekt.Provides(http.MethodGet, "/health",
	bidirekt.Response(http.StatusOK, nil),
)
