package fuego

// Serde implements serialization and deserialization for a given type.
type Serde interface {
	Serialize(v any) ([]byte, error)
	Deserialize(data []byte) (any, error)
}
