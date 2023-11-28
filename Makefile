default: ci

ci: fmt lint cover

ci-full: ci sec dependencies-analyze bench 

test: 
	go test ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

cover-web: cover
	go tool cover -html=coverage.out

bench:
	go test -bench ./... -benchmem

sec:
	go run github.com/securego/gosec/v2/cmd/gosec@latest ./...

dependencies-analyze:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

fmt:
	go run mvdan.cc/gofumpt@latest -l -w .

lint: 
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run

example:
	( cd examples/simple-crud && go run . -debug )

example-watch:
	( cd examples/simple-crud && air -- -debug )

# Documentation website
docs:
	go run golang.org/x/pkgsite/cmd/pkgsite@latest -http localhost:8084

docs-open:
	go run golang.org/x/pkgsite/cmd/pkgsite@latest -http localhost:8084 -open
