package internal

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommonContext_GetOpenAPIParams(t *testing.T) {
	params := map[string]OpenAPIParam{
		"test": {Name: "test", Description: "test param"},
	}
	ctx := CommonContext[any]{
		OpenAPIParams: params,
	}

	result := ctx.GetOpenAPIParams()
	require.Equal(t, params, result)
}

func TestCommonContext_Context(t *testing.T) {
	baseCtx := context.Background()
	ctx := CommonContext[any]{
		CommonCtx: baseCtx,
	}

	result := ctx.Context()
	require.Equal(t, baseCtx, result)
}

func TestCommonContext_Deadline(t *testing.T) {
	t.Run("context without deadline", func(t *testing.T) {
		ctx := CommonContext[any]{
			CommonCtx: context.Background(),
		}

		deadline, ok := ctx.Deadline()
		require.False(t, ok)
		require.True(t, deadline.IsZero())
	})

	t.Run("context with deadline", func(t *testing.T) {
		deadline := time.Now().Add(time.Hour)
		baseCtx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		ctx := CommonContext[any]{
			CommonCtx: baseCtx,
		}

		resultDeadline, ok := ctx.Deadline()
		require.True(t, ok)
		require.Equal(t, deadline, resultDeadline)
	})
}

func TestCommonContext_Done(t *testing.T) {
	t.Run("context not cancelled", func(t *testing.T) {
		ctx := CommonContext[any]{
			CommonCtx: context.Background(),
		}

		done := ctx.Done()
		require.Nil(t, done)
	})

	t.Run("context cancelled", func(t *testing.T) {
		baseCtx, cancel := context.WithCancel(context.Background())
		ctx := CommonContext[any]{
			CommonCtx: baseCtx,
		}

		done := ctx.Done()
		require.NotNil(t, done)

		cancel()
		select {
		case <-done:
			// Expected
		case <-time.After(time.Millisecond * 100):
			t.Fatal("context should be cancelled")
		}
	})
}

func TestCommonContext_Err(t *testing.T) {
	t.Run("context not cancelled", func(t *testing.T) {
		ctx := CommonContext[any]{
			CommonCtx: context.Background(),
		}

		err := ctx.Err()
		require.NoError(t, err)
	})

	t.Run("context cancelled", func(t *testing.T) {
		baseCtx, cancel := context.WithCancel(context.Background())
		ctx := CommonContext[any]{
			CommonCtx: baseCtx,
		}

		cancel()
		err := ctx.Err()
		require.Equal(t, context.Canceled, err)
	})
}

func TestCommonContext_Value(t *testing.T) {
	key := "test-key"
	value := "test-value"
	baseCtx := context.WithValue(context.Background(), key, value)
	ctx := CommonContext[any]{
		CommonCtx: baseCtx,
	}

	result := ctx.Value(key)
	require.Equal(t, value, result)

	result = ctx.Value("non-existent")
	require.Nil(t, result)
}

func TestCommonContext_QueryParams(t *testing.T) {
	urlValues := url.Values{
		"param1": []string{"value1"},
		"param2": []string{"value2a", "value2b"},
	}
	ctx := CommonContext[any]{
		UrlValues: urlValues,
	}

	result := ctx.QueryParams()
	require.Equal(t, urlValues, result)
}

func TestCommonContext_HasQueryParam(t *testing.T) {
	urlValues := url.Values{
		"existing": []string{"value"},
	}
	ctx := CommonContext[any]{
		UrlValues: urlValues,
	}

	require.True(t, ctx.HasQueryParam("existing"))
	require.False(t, ctx.HasQueryParam("non-existing"))
}

func TestCommonContext_QueryParam(t *testing.T) {
	t.Run("existing parameter", func(t *testing.T) {
		urlValues := url.Values{
			"test": []string{"value"},
		}
		ctx := CommonContext[any]{
			CommonCtx: context.Background(),
			UrlValues: urlValues,
			OpenAPIParams: map[string]OpenAPIParam{
				"test": {Name: "test"},
			},
		}

		result := ctx.QueryParam("test")
		require.Equal(t, "value", result)
	})

	t.Run("non-existing parameter with default", func(t *testing.T) {
		ctx := CommonContext[any]{
			CommonCtx: context.Background(),
			UrlValues: url.Values{},
			OpenAPIParams: map[string]OpenAPIParam{
				"test": {Name: "test", Default: "default-value"},
			},
		}

		result := ctx.QueryParam("test")
		require.Equal(t, "default-value", result)
	})

	t.Run("non-existing parameter without default", func(t *testing.T) {
		ctx := CommonContext[any]{
			CommonCtx: context.Background(),
			UrlValues: url.Values{},
			OpenAPIParams: map[string]OpenAPIParam{
				"test": {Name: "test"},
			},
		}

		result := ctx.QueryParam("test")
		require.Empty(t, result)
	})

	t.Run("parameter not in OpenAPI spec", func(t *testing.T) {
		urlValues := url.Values{
			"unexpected": []string{"value"},
		}
		ctx := CommonContext[any]{
			CommonCtx:     context.Background(),
			UrlValues:     urlValues,
			OpenAPIParams: map[string]OpenAPIParam{},
		}

		result := ctx.QueryParam("unexpected")
		require.Equal(t, "value", result)
	})
}

func TestCommonContext_QueryParamIntErr(t *testing.T) {
	t.Run("valid integer parameter", func(t *testing.T) {
		urlValues := url.Values{
			"count": []string{"42"},
		}
		ctx := CommonContext[any]{
			UrlValues: urlValues,
			OpenAPIParams: map[string]OpenAPIParam{
				"count": {Name: "count"},
			},
		}

		result, err := ctx.QueryParamIntErr("count")
		require.NoError(t, err)
		require.Equal(t, 42, result)
	})

	t.Run("invalid integer parameter", func(t *testing.T) {
		urlValues := url.Values{
			"count": []string{"invalid"},
		}
		ctx := CommonContext[any]{
			UrlValues: urlValues,
			OpenAPIParams: map[string]OpenAPIParam{
				"count": {Name: "count"},
			},
		}

		result, err := ctx.QueryParamIntErr("count")
		require.Error(t, err)
		require.Zero(t, result)

		var invalidErr QueryParamInvalidTypeError
		require.ErrorAs(t, err, &invalidErr)
		assert.Equal(t, "count", invalidErr.ParamName)
		assert.Equal(t, "invalid", invalidErr.ParamValue)
		assert.Equal(t, "int", invalidErr.ExpectedType)
	})

	t.Run("missing parameter with default", func(t *testing.T) {
		ctx := CommonContext[any]{
			UrlValues: url.Values{},
			OpenAPIParams: map[string]OpenAPIParam{
				"count": {Name: "count", Default: 10},
			},
		}

		result, err := ctx.QueryParamIntErr("count")
		require.NoError(t, err)
		require.Equal(t, 10, result)
	})

	t.Run("missing parameter without default", func(t *testing.T) {
		ctx := CommonContext[any]{
			UrlValues: url.Values{},
			OpenAPIParams: map[string]OpenAPIParam{
				"count": {Name: "count"},
			},
		}

		result, err := ctx.QueryParamIntErr("count")
		require.Error(t, err)
		require.Zero(t, result)

		var notFoundErr QueryParamNotFoundError
		require.ErrorAs(t, err, &notFoundErr)
		assert.Equal(t, "count", notFoundErr.ParamName)
	})
}

func TestCommonContext_QueryParamArr(t *testing.T) {
	t.Run("existing array parameter", func(t *testing.T) {
		urlValues := url.Values{
			"tags": []string{"tag1", "tag2", "tag3"},
		}
		ctx := CommonContext[any]{
			CommonCtx: context.Background(),
			UrlValues: urlValues,
			OpenAPIParams: map[string]OpenAPIParam{
				"tags": {Name: "tags"},
			},
		}

		result := ctx.QueryParamArr("tags")
		require.Equal(t, []string{"tag1", "tag2", "tag3"}, result)
	})

	t.Run("non-existing array parameter", func(t *testing.T) {
		ctx := CommonContext[any]{
			CommonCtx:     context.Background(),
			UrlValues:     url.Values{},
			OpenAPIParams: map[string]OpenAPIParam{},
		}

		result := ctx.QueryParamArr("tags")
		require.Nil(t, result)
	})

	t.Run("parameter not in OpenAPI spec", func(t *testing.T) {
		urlValues := url.Values{
			"unexpected": []string{"value1", "value2"},
		}
		ctx := CommonContext[any]{
			CommonCtx:     context.Background(),
			UrlValues:     urlValues,
			OpenAPIParams: map[string]OpenAPIParam{},
		}

		result := ctx.QueryParamArr("unexpected")
		require.Equal(t, []string{"value1", "value2"}, result)
	})
}

func TestCommonContext_QueryParamInt(t *testing.T) {
	t.Run("valid integer parameter", func(t *testing.T) {
		urlValues := url.Values{
			"count": []string{"42"},
		}
		ctx := CommonContext[any]{
			UrlValues: urlValues,
			OpenAPIParams: map[string]OpenAPIParam{
				"count": {Name: "count"},
			},
		}

		result := ctx.QueryParamInt("count")
		require.Equal(t, 42, result)
	})

	t.Run("invalid integer parameter returns 0", func(t *testing.T) {
		urlValues := url.Values{
			"count": []string{"invalid"},
		}
		ctx := CommonContext[any]{
			UrlValues: urlValues,
			OpenAPIParams: map[string]OpenAPIParam{
				"count": {Name: "count"},
			},
		}

		result := ctx.QueryParamInt("count")
		require.Zero(t, result)
	})
}

func TestCommonContext_QueryParamBoolErr(t *testing.T) {
	t.Run("valid boolean parameter", func(t *testing.T) {
		urlValues := url.Values{
			"active": []string{"true"},
		}
		ctx := CommonContext[any]{
			UrlValues: urlValues,
			OpenAPIParams: map[string]OpenAPIParam{
				"active": {Name: "active"},
			},
		}

		result, err := ctx.QueryParamBoolErr("active")
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("invalid boolean parameter", func(t *testing.T) {
		urlValues := url.Values{
			"active": []string{"invalid"},
		}
		ctx := CommonContext[any]{
			UrlValues: urlValues,
			OpenAPIParams: map[string]OpenAPIParam{
				"active": {Name: "active"},
			},
		}

		result, err := ctx.QueryParamBoolErr("active")
		require.Error(t, err)
		require.False(t, result)

		var invalidErr QueryParamInvalidTypeError
		require.ErrorAs(t, err, &invalidErr)
		assert.Equal(t, "active", invalidErr.ParamName)
		assert.Equal(t, "invalid", invalidErr.ParamValue)
		assert.Equal(t, "bool", invalidErr.ExpectedType)
	})

	t.Run("missing parameter with default", func(t *testing.T) {
		ctx := CommonContext[any]{
			UrlValues: url.Values{},
			OpenAPIParams: map[string]OpenAPIParam{
				"active": {Name: "active", Default: true},
			},
		}

		result, err := ctx.QueryParamBoolErr("active")
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("missing parameter without default", func(t *testing.T) {
		ctx := CommonContext[any]{
			UrlValues: url.Values{},
			OpenAPIParams: map[string]OpenAPIParam{
				"active": {Name: "active"},
			},
		}

		result, err := ctx.QueryParamBoolErr("active")
		require.Error(t, err)
		require.False(t, result)

		var notFoundErr QueryParamNotFoundError
		require.ErrorAs(t, err, &notFoundErr)
		assert.Equal(t, "active", notFoundErr.ParamName)
	})
}

func TestCommonContext_QueryParamBool(t *testing.T) {
	t.Run("valid boolean parameter", func(t *testing.T) {
		urlValues := url.Values{
			"active": []string{"true"},
		}
		ctx := CommonContext[any]{
			UrlValues: urlValues,
			OpenAPIParams: map[string]OpenAPIParam{
				"active": {Name: "active"},
			},
		}

		result := ctx.QueryParamBool("active")
		require.True(t, result)
	})

	t.Run("invalid boolean parameter returns false", func(t *testing.T) {
		urlValues := url.Values{
			"active": []string{"invalid"},
		}
		ctx := CommonContext[any]{
			UrlValues: urlValues,
			OpenAPIParams: map[string]OpenAPIParam{
				"active": {Name: "active"},
			},
		}

		result := ctx.QueryParamBool("active")
		require.False(t, result)
	})
}

func TestQueryParamNotFoundError(t *testing.T) {
	err := QueryParamNotFoundError{ParamName: "test"}

	assert.Equal(t, "param test not found", err.Error())
	assert.Equal(t, 422, err.StatusCode())
	assert.Equal(t, "param test not found", err.DetailMsg())
}

func TestQueryParamInvalidTypeError(t *testing.T) {
	innerErr := assert.AnError
	err := QueryParamInvalidTypeError{
		ParamName:    "count",
		ParamValue:   "invalid",
		ExpectedType: "int",
		Err:          innerErr,
	}

	assert.Equal(t, "query param count=invalid is not of type int: assert.AnError general error for testing", err.Error())
	assert.Equal(t, 422, err.StatusCode())
	assert.Equal(t, "query param count=invalid is not of type int", err.DetailMsg())
}
