package fuego

import (
	"context"
	"io"
)

// Serde implements serialization and deserialization for a given content type.
type Serde interface {
	Serialize(v any) ([]byte, error)
	Deserialize(ctx context.Context, input io.Reader) (any, error)
}
