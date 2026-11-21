module example.com/api

go 1.26.0

require (
	github.com/bidirekt/go-sdk v0.0.0
	example.com/shared v0.0.0
)

replace (
	github.com/bidirekt/go-sdk => ../../../../../..
	example.com/shared => ../shared
)
