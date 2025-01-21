package data

import (
	_ "embed"
)

// Favicon is the embedded SVG favicon file.
// Will not pollute your binary with a favicon.ico file,
// unless you import it explicitly in your server: the Fuego framework does not use it.
//
//go:embed fuego.svg
var Favicon []byte
