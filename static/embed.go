package data

import (
	_ "embed"
)

// Embeded SVG favicon file.
// Will not pollute your binary with a favicon.ico file,
// unless you import it explicitely in your server: the Fuego framework does not use it.
//
//go:embed fuego.svg
var Favicon []byte
