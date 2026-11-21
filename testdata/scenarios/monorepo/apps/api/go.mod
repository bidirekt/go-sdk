module example.com/acme/apps/api

go 1.26.0

require (
	example.com/acme/shared v0.0.0
	github.com/bidirekt/go-sdk v0.0.0-00010101000000-000000000000
)

require (
	github.com/bidirekt/go-sdk/cmd/go-sdk v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/mod v0.43.0 // indirect
	golang.org/x/sync v0.25.0 // indirect
	golang.org/x/tools v0.53.0 // indirect
)

replace example.com/acme/shared => ../../shared

replace github.com/bidirekt/go-sdk => ../../../../..

replace github.com/bidirekt/go-sdk/cmd/go-sdk => ../../../../../cmd/go-sdk

tool github.com/bidirekt/go-sdk/cmd/go-sdk
