test: 
	go test

cover:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
	rm coverage.out

fuzz:
	go test -fuzz Fuzz -fuzztime 10s

fmt:
	gofumpt -l -w .

example:
	go run ./examples/simple-crud/main.go
