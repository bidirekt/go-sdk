package missinggosum

import (
	"github.com/bidirekt/go-sdk"
	"golang.org/x/sync/errgroup"
)

var group errgroup.Group

var _ = bidirekt.Provides("GET", "/pets",
	bidirekt.Response(200, nil),
)
