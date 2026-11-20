package petstore

import (
	"example.com/audit"
	"github.com/bidirekt/go-sdk"
)

type Timestamps struct {
	CreatedAt string `json:"createdAt" bidirekt:"field"`
}

type Pet struct {
	Timestamps
	Name    string           `json:"name" bidirekt:"field"`
	History audit.Timestamps `json:"history" bidirekt:"field"`
}

var _ = bidirekt.Provides("GET", "/pets/*",
	bidirekt.Response(200, Pet{}),
)
