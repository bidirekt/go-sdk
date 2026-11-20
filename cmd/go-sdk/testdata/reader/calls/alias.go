package calls

import (
	"net/http"

	sdk "github.com/bidirekt/go-sdk"
)

var _ = sdk.Provides(http.MethodGet, "/alias",
	sdk.Response(http.StatusOK, nil),
)
