package calls

import (
	. "github.com/bidirekt/go-sdk"
)

var _ = Consumes("petstore_reviews", "POST", "/dot",
	Request(nil),
	Response(201, nil),
)
