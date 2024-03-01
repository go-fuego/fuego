mods=$(go list -f '{{.Dir}}' -m | xargs)
for mod in $mods; do
    (cd "$mod";  go get -u all ; go mod tidy ; go test ./...)
done
