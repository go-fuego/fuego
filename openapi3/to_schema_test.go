package openapi3

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestToSchema(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		s := ToSchema("")
		require.Equal(t, String, s.Type)
	})

	t.Run("alias to string", func(t *testing.T) {
		type S string
		s := ToSchema(S(""))
		require.Equal(t, String, s.Type)
	})

	t.Run("struct with a field alias to string", func(t *testing.T) {
		type MyAlias string
		type S struct {
			A MyAlias
		}

		s := ToSchema(S{})
		require.Equal(t, Object, s.Type)
		require.Equal(t, String, s.Properties["A"].Type)
	})

	t.Run("int", func(t *testing.T) {
		s := ToSchema(0)
		require.Equal(t, Integer, s.Type)
	})

	t.Run("bool", func(t *testing.T) {
		s := ToSchema(false)
		require.Equal(t, Boolean, s.Type)
	})

	t.Run("time", func(t *testing.T) {
		s := ToSchema(time.Now())
		require.Equal(t, String, s.Type)
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
		require.Equal(t, Object, s.Type)
		require.Equal(t, String, s.Properties["a"].Type)
		require.Equal(t, Integer, s.Properties["B"].Type)
		require.Equal(t, Boolean, s.Properties["C"].Type)
		require.Equal(t, []string{"a"}, s.Required)
		require.Equal(t, Object, s.Properties["Nested"].Type)
		require.Equal(t, String, s.Properties["Nested"].Properties["C"].Type)

		gotSchema, err := json.Marshal(s)
		require.NoError(t, err)
		require.JSONEq(t, string(gotSchema), `
		{
			"type":"object",
			"required":["a"],
			"properties": {
				"a":{"type":"string","examples":["hello"]},
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
		require.Equal(t, Object, s.Type)
		require.Equal(t, String, s.Properties["A"].Type)
		require.Equal(t, Integer, s.Properties["B"].Type)
		// TODO require.Equal(t, []string{"A", "B", "Nested"}, s.Required)
		require.Equal(t, Object, s.Properties["Nested"].Type)
		require.Equal(t, String, s.Properties["Nested"].Properties["C"].Type)

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

	t.Run("slice of strings", func(t *testing.T) {
		s := ToSchema([]string{})
		require.Equal(t, Array, s.Type)
		require.Equal(t, String, s.Items.Type)
	})

	t.Run("slice of structs", func(t *testing.T) {
		type S struct {
			A string
		}
		s := ToSchema([]S{})
		require.Equal(t, Array, s.Type)
		require.Equal(t, Object, s.Items.Type)
		require.Equal(t, String, s.Items.Properties["A"].Type)
	})

	t.Run("slice of ptr to struct", func(t *testing.T) {
		type S struct {
			A string
		}
		s := ToSchema([]*S{})
		require.Equal(t, Array, s.Type)
		require.Equal(t, Object, s.Items.Type)
		require.Equal(t, String, s.Items.Properties["A"].Type)
	})

	t.Run("embedded struct", func(t *testing.T) {
		type S struct {
			A string
		}
		type T struct {
			S
			B int
		}
		s := ToSchema(T{})
		require.Equal(t, Object, s.Type)
		require.Equal(t, OpenAPIType(""), s.Properties["A"].Type)
		require.Equal(t, Object, s.Properties["S"].Type)
		require.Equal(t, String, s.Properties["S"].Properties["A"].Type)
		require.Equal(t, Integer, s.Properties["B"].Type)
	})

	t.Run("struct of slices of structs", func(t *testing.T) {
		type S struct {
			A string
		}

		type T struct {
			SliceOfS []S
		}

		tt := ToSchema(T{})
		require.Equal(t, Object, tt.Type)
		require.Equal(t, Array, tt.Properties["SliceOfS"].Type)
		require.NotNil(t, tt.Properties["SliceOfS"].Items)
		require.Equal(t, Array, tt.Properties["SliceOfS"].Type)
		require.Equal(t, Object, tt.Properties["SliceOfS"].Items.Type)
		require.Equal(t, String, tt.Properties["SliceOfS"].Items.Properties["A"].Type)
	})

	t.Run("struct with ptrs properties", func(t *testing.T) {
		t.Skip("TODO")
		type S struct {
			A *string
			B *int
		}
		s := ToSchema(S{})
		require.Equal(t, Object, s.Type)
		require.Equal(t, String, s.Properties["A"].Type)
		require.Equal(t, Integer, s.Properties["B"].Type)
	})
}

func TestFieldName(t *testing.T) {
	t.Run("no tag", func(t *testing.T) {
		require.Equal(t, "A", fieldName(reflect.StructField{
			Name: "A",
		}))
	})

	t.Run("json tag", func(t *testing.T) {
		require.Equal(t, "a", fieldName(reflect.StructField{
			Name: "A",
			Tag:  `json:"a"`,
		}))
	})

	t.Run("json tag with omitempty", func(t *testing.T) {
		require.Equal(t, "a", fieldName(reflect.StructField{
			Name: "A",
			Tag:  `json:"a,omitempty"`,
		}))
	})

	t.Run("json tag with no name, omitempty", func(t *testing.T) {
		require.Equal(t, "A", fieldName(reflect.StructField{
			Name: "A",
			Tag:  `json:",omitempty"`,
		}))
	})
}
