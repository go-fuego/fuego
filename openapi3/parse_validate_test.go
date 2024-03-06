package openapi3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseValidate(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := &Schema{}
		parseValidate(s, "")
	})

	t.Run("required", func(t *testing.T) {
		t.Skip("TODO")
		s := &Schema{}
		parseValidate(s, "required")
		require.Equal(t, true, s.Required)
	})

	t.Run("min for int", func(t *testing.T) {
		s := &Schema{
			Type: "integer",
		}
		parseValidate(s, "min=10")
		require.Equal(t, 10, s.Minimum)
	})

	t.Run("max for int", func(t *testing.T) {
		s := &Schema{
			Type: "integer",
		}
		parseValidate(s, "max=10")
		require.Equal(t, 10, s.Maximum)
	})

	t.Run("min for string", func(t *testing.T) {
		s := &Schema{
			Type: "string",
		}
		parseValidate(s, "min=10")
		require.Equal(t, 10, s.MinLength)
	})

	t.Run("max for string", func(t *testing.T) {
		s := &Schema{
			Type: "string",
		}
		parseValidate(s, "max=10")
		require.Equal(t, 10, s.MaxLength)
	})

	t.Run("multiple", func(t *testing.T) {
		s := &Schema{
			Type: "string",
		}
		parseValidate(s, "required,min=10,max=20")
		require.Equal(t, 10, s.MinLength)
		require.Equal(t, 20, s.MaxLength)
	})
}
