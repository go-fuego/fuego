package fuego

import (
	"context"
	"io"
)

// SerDes implements serialization and deserialization for a given content type.
// A round trip of SerDes is expected to be idempotent.
// This means that the result of Serialize(Deserialize(v)) should be equal to v.
type SerDes interface {

	// Serialize serializes the value v into a byte slice.
	// The byte slice is expected to be a valid representation of the value v for the content-type.
	// If the value v is not valid, an error is returned.
	Serialize(v any) ([]byte, error)

	// Deserialize deserializes the input into a value of type T.
	// The input is expected to be a valid representation of the type T for the content-type.
	// If the input is not valid, an error is returned.
	Deserialize(ctx context.Context, input io.Reader) (any, error)
}
