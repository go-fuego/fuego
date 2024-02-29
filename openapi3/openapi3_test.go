package openapi3

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToSchema(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		s := ToSchema("")
		if s.Type != "string" {
			t.Errorf("expected string, got %s", s.Type)
		}
	})

	t.Run("int", func(t *testing.T) {
		s := ToSchema(0)
		if s.Type != "integer" {
			t.Errorf("expected integer, got %s", s.Type)
		}
	})

	t.Run("struct", func(t *testing.T) {
		type S struct {
			A      string `json:"a" validate:"required" example:"hello"`
			B      int
			C      bool
			Nested struct {
				C string
			}
		}
		s := ToSchema(S{})
		require.Equal(t, "object", s.Type)
		require.Equal(t, "string", s.Properties["a"].Type)
		require.Equal(t, "integer", s.Properties["B"].Type)
		require.Equal(t, "boolean", s.Properties["C"].Type)
		require.Equal(t, []string{"a"}, s.Required)
		require.Equal(t, "object", s.Properties["Nested"].Type)
		require.Equal(t, "string", s.Properties["Nested"].Properties["C"].Type)

		gotSchema, err := json.Marshal(s)
		require.NoError(t, err)
		require.JSONEq(t, string(gotSchema), `
		{
			"type":"object",
			"required":["a"],
			"properties": {
				"a":{"type":"string","example":"hello"},
				"B":{"type":"integer"},
				"C":{"type":"boolean"},
				"Nested":{
					"type":"object",
					"properties":{
						"C":{"type":"string"}
					}
				}
			}
		}`)
	})

	t.Run("ptr to struct", func(t *testing.T) {
		type S struct {
			A      string
			B      int
			Nested struct {
				C string
			}
		}
		s := ToSchema(&S{})
		require.Equal(t, "object", s.Type)
		require.Equal(t, "string", s.Properties["A"].Type)
		require.Equal(t, "integer", s.Properties["B"].Type)
		// TODO require.Equal(t, []string{"A", "B", "Nested"}, s.Required)
		require.Equal(t, "object", s.Properties["Nested"].Type)
		require.Equal(t, "string", s.Properties["Nested"].Properties["C"].Type)

		gotSchema, err := json.Marshal(s)
		require.NoError(t, err)
		require.JSONEq(t, string(gotSchema), `
		{
			"type":"object",
			"properties": {
				"A":{"type":"string"},
				"B":{"type":"integer"},
				"Nested":{
					"type":"object",
					"properties":{
						"C":{"type":"string"}
					}
				}
			}
		}`)
	})
}
