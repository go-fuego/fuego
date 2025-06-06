package fuego

import (
	"context"
	"io"
)

// SerDes implements serialization and deserialization for a given content type.
type SerDes interface {
	Serialize(v any) ([]byte, error)
	Deserialize(ctx context.Context, input io.Reader) (any, error)
}
