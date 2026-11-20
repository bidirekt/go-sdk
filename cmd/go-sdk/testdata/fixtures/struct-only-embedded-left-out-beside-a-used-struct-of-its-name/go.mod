module example.com/petstore

go 1.26.0

require (
	github.com/bidirekt/go-sdk v0.0.0
	example.com/audit v0.0.0
)

replace (
	github.com/bidirekt/go-sdk => ../../../../..
	example.com/audit => ./audit
)
