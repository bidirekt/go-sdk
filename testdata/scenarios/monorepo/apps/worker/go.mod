module example.com/acme/apps/worker

go 1.24.0

require example.com/acme/shared v0.0.0

require github.com/bidirekt/go-sdk v0.0.0-00010101000000-000000000000 // indirect

replace example.com/acme/shared => ../../shared

replace github.com/bidirekt/go-sdk => ../../../../..
