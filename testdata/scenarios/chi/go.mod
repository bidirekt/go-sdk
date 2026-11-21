module example.com/bookshelf

go 1.26.0

replace github.com/bidirekt/go-sdk => ../../..

replace github.com/bidirekt/go-sdk/cmd/go-sdk => ../../../cmd/go-sdk

require (
	github.com/bidirekt/go-sdk v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.5.0
)

require (
	github.com/bidirekt/go-sdk/cmd/go-sdk v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/mod v0.43.0 // indirect
	golang.org/x/sync v0.25.0 // indirect
	golang.org/x/tools v0.53.0 // indirect
)

tool github.com/bidirekt/go-sdk/cmd/go-sdk
