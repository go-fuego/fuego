#!/bin/bash

mods=$(go list -f '{{.Dir}}' -m)
for mod in $mods; do
    go test -C "$mod" ./...
done
