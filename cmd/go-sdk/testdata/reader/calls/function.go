package calls

import (
	"github.com/bidirekt/go-sdk"
)

func routes() {
	_ = bidirekt.Provides("PUT", "/inside",
		bidirekt.Response(204, nil),
	)
}
