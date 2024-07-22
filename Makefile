default: ci

ci: fmt lint cover

ci-full: ci dependencies-analyze openapi-check check-all-modules lint-markdown bench

test: 
	go test ./...

check-all-modules:
	./check-all-modules.sh

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	sed '/petstore/d' coverage.out > coverage.out

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
	markdownlint --ignore documentation/node_modules --dot .

# Update golden files
golden-update:
	(cd examples/petstore && go test -update)

# Check OpenAPI spec generated for the Petstore example. Uses https://github.com/daveshanley/vacuum
openapi-check:
	vacuum lint -d examples/petstore/testdata/doc/openapi.json

# Examples
example:
	( cd examples/full-app-gourmet && go run . -debug )

example-watch:
	( cd examples/full-app-gourmet && air -- -debug )

petstore:
	( cd examples/petstore && go run . -debug )

# Documentation website
docs:
	go run golang.org/x/pkgsite/cmd/pkgsite@latest -http localhost:8084

docs-open:
	go run golang.org/x/pkgsite/cmd/pkgsite@latest -http localhost:8084 -open

.PHONY: docs-open docs example-watch example lint lint-markdown fmt ci ci-full
.PHONY: dependencies-analyze build bench cover-web cover test petstore check-all-modules
.PHONY: golden-update openapi-check
