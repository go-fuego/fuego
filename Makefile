default: ci

ci: fmt lint cover

ci-full: ci dependencies-analyze bench

test: 
	go test ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

cover-web: cover
	go tool cover -html=coverage.out

bench:
	go test -bench ./... -benchmem

build:
	go build -v ./...

dependencies-analyze:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

fmt:
	go run mvdan.cc/gofumpt@latest -l -w -extra .

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run

lint-markdown:
	markdownlint --dot .

example:
	( cd examples/full-app-gourmet && go run . -debug )

example-watch:
	( cd examples/full-app-gourmet && air -- -debug )

# Documentation website
docs:
	go run golang.org/x/pkgsite/cmd/pkgsite@latest -http localhost:8084

docs-open:
	go run golang.org/x/pkgsite/cmd/pkgsite@latest -http localhost:8084 -open

.PHONY: docs-open docs example-watch example lint lint-markdown fmt ci ci-full
.PHONY: dependencies-analyze build bench cover-web cover test
