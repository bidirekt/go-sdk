module example.com/shop

go 1.26.0

replace github.com/bidirekt/go-sdk => ../../..

replace github.com/bidirekt/go-sdk/cmd/go-sdk => ../../../cmd/go-sdk

require (
	github.com/bidirekt/go-sdk v0.0.0-00010101000000-000000000000
	github.com/labstack/echo/v4 v4.16.6
)

require (
	github.com/bidirekt/go-sdk/cmd/go-sdk v0.0.0-00010101000000-000000000000 // indirect
	github.com/labstack/gommon v0.5.0 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.22 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.59.0 // indirect
	golang.org/x/mod v0.43.0 // indirect
	golang.org/x/net v0.63.0 // indirect
	golang.org/x/sync v0.25.0 // indirect
	golang.org/x/sys v0.50.0 // indirect
	golang.org/x/text v0.44.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	golang.org/x/tools v0.53.0 // indirect
)

tool github.com/bidirekt/go-sdk/cmd/go-sdk
