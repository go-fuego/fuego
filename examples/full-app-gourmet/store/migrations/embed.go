package migrations

import (
	"embed"
)

//go:embed *.sql
var FS embed.FS
