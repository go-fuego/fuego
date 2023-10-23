default: ci

ci: fmt lint cover

ci-full: ci bench fuzz

test: 
	go test

cover:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out

cover-web: cover
	go tool cover -html=coverage.out

bench:
	go test -bench . -benchmem

fuzz:
	go test -fuzz Fuzz -fuzztime 10s

fmt:
	go run mvdan.cc/gofumpt@latest -l -w .

# If golangci-lint is not installed, run it from latest github version found. Installed version is faster.
lint: 
	golangci-lint run || (echo "Running golangci-lint from latest github version found" && go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run)

example:
	go run ./examples/simple-crud
