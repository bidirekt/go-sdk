package calls

import (
	"example.com/calls/fake"
)

var _ = bidirekt.Provides("GET", "/foreign")
