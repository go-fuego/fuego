package fuego

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RegisterParameters(t *testing.T) {
	t.Run("Add parameters to empty operation", func(t *testing.T) {
		operation := &openapi3.Operation{}

		param1 := &openapi3.Parameter{
			Name:   "testParam1",
			In:     "query",
			Schema: &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
		}

		err := RegisterParameters(operation, param1)

		require.NoError(t, err)
		assert.Len(t, operation.Parameters, 1)
		assert.Equal(t, param1, operation.Parameters[0].Value)
	})

	t.Run("Add multiple parameters", func(t *testing.T) {
		operation := &openapi3.Operation{}

		param1 := &openapi3.Parameter{
			Name:   "testParam1",
			In:     "query",
			Schema: &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
		}

		param2 := &openapi3.Parameter{
			Name:   "testParam2",
			In:     "header",
			Schema: &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
		}

		err := RegisterParameters(operation, param1, param2)

		require.NoError(t, err)
		assert.Len(t, operation.Parameters, 2)
		assert.Equal(t, param1, operation.Parameters[0].Value)
		assert.Equal(t, param2, operation.Parameters[1].Value)
	})

	t.Run("Add parameters to operation with existing parameters", func(t *testing.T) {
		operation := &openapi3.Operation{
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:   "existingParam",
						In:     "path",
						Schema: &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
					},
				},
			},
		}

		param1 := &openapi3.Parameter{
			Name:   "testParam1",
			In:     "query",
			Schema: &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
		}

		err := RegisterParameters(operation, param1)

		require.NoError(t, err)
		assert.Len(t, operation.Parameters, 2)
		assert.Equal(t, "existingParam", operation.Parameters[0].Value.Name)
		assert.Equal(t, "testParam1", operation.Parameters[1].Value.Name)
	})

	t.Run("Nil parameter results in error", func(t *testing.T) {
		operation := &openapi3.Operation{}

		err := RegisterParameters(operation, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parameter cannot be nil")
		assert.Len(t, operation.Parameters, 0)
	})

	t.Run("Nil parameter in multi-parameter call", func(t *testing.T) {
		operation := &openapi3.Operation{}

		param1 := &openapi3.Parameter{
			Name:   "testParam1",
			In:     "query",
			Schema: &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
		}

		err := RegisterParameters(operation, param1, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parameter cannot be nil")
		assert.Len(t, operation.Parameters, 0)
	})
}
