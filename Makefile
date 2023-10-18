default: ci

ci: fmt lint cover

ci-full: ci fuzz

test: 
	go test

cover:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
	rm coverage.out

bench:
	go test -bench . -benchmem

fuzz:
	go test -fuzz Fuzz -fuzztime 10s

fmt:
	gofumpt -l -w .

lint:
	golangci-lint run

example:
	go run ./examples/simple-crud/main.go
