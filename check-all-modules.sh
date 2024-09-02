#!/bin/bash

set -euo pipefail

mods=$(go list -f '{{.Dir}}' -m)
for mod in $mods; do
	cd "$mod"
	echo "=== Updating $mod"
    go get -u ./...
	go mod tidy
	go test ./...
	go build -o /dev/null ./...
	cd -
	echo
done
