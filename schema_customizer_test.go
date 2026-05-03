package fuego

import (
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetermineFieldConstraints(t *testing.T) {
	t.Run("non-struct type is a no-op", func(t *testing.T) {
		schema := &openapi3.Schema{}
		determineFieldConstraints(reflect.TypeFor[string](), schema)
		assert.Empty(t, schema.Required)
	})
	t.Run("private field is skipped", func(t *testing.T) {
		type S struct {
			private string //nolint:unused
		}
		schema := &openapi3.Schema{
			Properties: openapi3.Schemas{},
		}
		determineFieldConstraints(reflect.TypeFor[S](), schema)
		assert.Empty(t, schema.Required)
	})
	t.Run("json tag - is skipped", func(t *testing.T) {
		type S struct {
			Hidden string `json:"-"`
		}
		schema := &openapi3.Schema{
			Properties: openapi3.Schemas{},
		}
		determineFieldConstraints(reflect.TypeFor[S](), schema)
		assert.Empty(t, schema.Required)
	})

	t.Run("field without omitempty is required", func(t *testing.T) {
		type S struct {
			Name string `json:"name"`
		}
		schema := &openapi3.Schema{
			Properties: openapi3.Schemas{
				"name": &openapi3.SchemaRef{Value: &openapi3.Schema{}},
			},
		}
		determineFieldConstraints(reflect.TypeFor[S](), schema)
		assert.Contains(t, schema.Required, "name")
	})

	t.Run("field with omitempty is not required", func(t *testing.T) {
		type S struct {
			Name string `json:"name,omitempty"`
		}
		schema := &openapi3.Schema{
			Properties: openapi3.Schemas{
				"name": &openapi3.SchemaRef{Value: &openapi3.Schema{}},
			},
		}
		determineFieldConstraints(reflect.TypeFor[S](), schema)
		assert.NotContains(t, schema.Required, "name")
	})

	t.Run("field with validate:\"required\" is required", func(t *testing.T) {
		type S struct {
			Name string `json:"name,omitempty" validate:"required"`
		}
		schema := &openapi3.Schema{
			Properties: openapi3.Schemas{
				"name": &openapi3.SchemaRef{Value: &openapi3.Schema{}},
			},
		}
		determineFieldConstraints(reflect.TypeFor[S](), schema)
		assert.Contains(t, schema.Required, "name")
		assert.False(t, schema.Properties["name"].Value.Nullable)
	})

	t.Run("slice field is nullable", func(t *testing.T) {
		type S struct {
			Items []string `json:"items"`
		}
		schema := &openapi3.Schema{
			Properties: openapi3.Schemas{
				"items": &openapi3.SchemaRef{Value: &openapi3.Schema{}},
			},
		}
		determineFieldConstraints(reflect.TypeFor[S](), schema)
		assert.True(t, schema.Properties["items"].Value.Nullable)
	})

	t.Run("map field is nullable", func(t *testing.T) {
		type S struct {
			Meta map[string]string `json:"meta"`
		}
		schema := &openapi3.Schema{
			Properties: openapi3.Schemas{
				"meta": &openapi3.SchemaRef{Value: &openapi3.Schema{}},
			},
		}
		determineFieldConstraints(reflect.TypeFor[S](), schema)
		assert.True(t, schema.Properties["meta"].Value.Nullable)
	})

	t.Run("string field is not nullable", func(t *testing.T) {
		type S struct {
			Name string `json:"name"`
		}
		schema := &openapi3.Schema{
			Properties: openapi3.Schemas{
				"name": &openapi3.SchemaRef{Value: &openapi3.Schema{}},
			},
		}
		determineFieldConstraints(reflect.TypeFor[S](), schema)
		assert.False(t, schema.Properties["name"].Value.Nullable)
	})

	t.Run("required fields are sorted", func(t *testing.T) {
		type S struct {
			Zebra string `json:"zebra"`
			Apple string `json:"apple"`
			Mango string `json:"mango"`
		}
		schema := &openapi3.Schema{
			Properties: openapi3.Schemas{
				"zebra": &openapi3.SchemaRef{Value: &openapi3.Schema{}},
				"apple": &openapi3.SchemaRef{Value: &openapi3.Schema{}},
				"mango": &openapi3.SchemaRef{Value: &openapi3.Schema{}},
			},
		}
		determineFieldConstraints(reflect.TypeFor[S](), schema)
		require.Len(t, schema.Required, 3)
		assert.Equal(t, []string{"apple", "mango", "zebra"}, schema.Required)
	})
}

func TestParseValidateOneof(t *testing.T) {
	t.Run("string oneof adds enum values", func(t *testing.T) {
		schema := &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}}
		parseValidate(reflect.StructTag(`validate:"oneof=a b c"`), schema)
		assert.Equal(t, []any{"a", "b", "c"}, schema.Enum)
	})

	t.Run("integer oneof adds typed enum values", func(t *testing.T) {
		schema := &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeInteger}}
		parseValidate(reflect.StructTag(`validate:"oneof=1 2 3"`), schema)
		assert.Equal(t, []any{1, 2, 3}, schema.Enum)
	})

	t.Run("number oneof adds typed enum values", func(t *testing.T) {
		schema := &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeNumber}}
		parseValidate(reflect.StructTag(`validate:"oneof=1.5 2.5"`), schema)
		assert.Equal(t, []any{1.5, 2.5}, schema.Enum)
	})

	t.Run("oneof works alongside required", func(t *testing.T) {
		schema := &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}}
		parseValidate(reflect.StructTag(`validate:"required,oneof=foo bar"`), schema)
		assert.Equal(t, []any{"foo", "bar"}, schema.Enum)
	})

	t.Run("no oneof leaves enum nil", func(t *testing.T) {
		schema := &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}}
		parseValidate(reflect.StructTag(`validate:"required"`), schema)
		assert.Nil(t, schema.Enum)
	})
}
