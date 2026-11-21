module example.com/app

go 1.26.0

require (
	example.com/money v0.0.0
	github.com/bidirekt/go-sdk v0.0.0
)

replace (
	example.com/money => ./money
	github.com/bidirekt/go-sdk => ../../../../..
)
