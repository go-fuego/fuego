package store

import "encoding/json"

// Override User's JSON serialization to exclude the encrypted password.
func (u User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		Alias
		EncryptedPassword string `json:"encrypted_password,omitempty"`
	}{
		Alias:             (Alias)(u),
		EncryptedPassword: "",
	})
}
