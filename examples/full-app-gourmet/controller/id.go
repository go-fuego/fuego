package controller

import (
	"crypto/rand"
	"encoding/hex"
)

func generateID() string {
	id := make([]byte, 10)
	_, _ = rand.Read(id)
	return hex.EncodeToString(id)
}
