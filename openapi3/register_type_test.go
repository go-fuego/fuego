package openapi3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocument_RegisterType(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		type S struct {
			A string
		}
		type T struct {
			S
			B int
		}

		d := NewDocument()
		s := d.RegisterType(T{})
		// will return a schema with a reference to the schema of T
		require.Equal(t, "#/components/schemas/T", s.Ref)
		require.Equal(t, "object", d.Components.Schemas["T"].Type)
		require.Equal(t, "integer", d.Components.Schemas["T"].Properties["B"].Type)
		require.Equal(t, "", d.Components.Schemas["T"].Properties["A"].Type)
		require.Equal(t, "string", d.Components.Schemas["T"].Properties["S"].Properties["A"].Type)
	})

	t.Run("array", func(t *testing.T) {
		type S struct {
			A string
		}

		d := NewDocument()
		s := d.RegisterType([]S{})
		// will return a schema with a reference to the schema of T
		require.Equal(t, "array", s.Type)
		require.Equal(t, "#/components/schemas/S", s.Items.Ref)
		require.Equal(t, "object", d.Components.Schemas["S"].Type)
		require.Equal(t, "string", d.Components.Schemas["S"].Properties["A"].Type)
	})
}

func BenchmarkDocument_RegisterType(b *testing.B) {
	d := NewDocument()
	type S struct {
		A string
	}
	for range b.N {
		d.RegisterType(S{})
	}
}
