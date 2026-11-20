package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type Health struct {
	Status string `json:"status" bidirekt:"field"`
}

var _ = bidirekt.Provides("GET", "/health",
	bidirekt.Response(200, struct {
		Status string `json:"status" bidirekt:"field"`
	}{}, bidirekt.Named("HealthReport")),
)

var _ = bidirekt.Provides("GET", "/ready",
	bidirekt.Response(200, struct {
		Status string `json:"status" bidirekt:"field"`
	}{}, bidirekt.Named("HealthReport")),
)
