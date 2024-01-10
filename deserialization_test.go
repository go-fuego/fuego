package fuego

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type BodyTest struct {
	A string
	B int
	C bool
}

func TestReadJSON(t *testing.T) {
	input := strings.NewReader(`{"A":"a","B":1,"C":true}`)

	t.Run("ReadJSON", func(t *testing.T) {
		body, err := ReadJSON[BodyTest](input)
		require.NoError(t, err)
		require.Equal(t, BodyTest{"a", 1, true}, body)
	})

	t.Run("cannot read invalid JSON", func(t *testing.T) {
		_, err := ReadJSON[BodyTest](input)
		require.ErrorAs(t, err, &BadRequestError{}, "Expected a BadRequestError")
	})

	t.Run("cannot deserialize JSON to wrong struct", func(t *testing.T) {
		type WrongBody struct {
			A string
			B int
			// Missing C bool
		}
		_, err := ReadJSON[WrongBody](input)
		require.ErrorAs(t, err, &BadRequestError{}, "Expected a BadRequestError")
	})
}

func TestReadString(t *testing.T) {
	t.Run("read string", func(t *testing.T) {
		input := strings.NewReader(`string decoded as is`)
		_, err := ReadString[string](input)
		require.NoError(t, err)
	})

	t.Run("read string alias", func(t *testing.T) {
		type StringAlias string
		input := strings.NewReader(`string decoded as is`)
		_, err := ReadString[StringAlias](input)
		require.NoError(t, err)
	})
}

func BenchmarkReadJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		input := strings.NewReader(`{"A":"a","B":1,"C":true}`)
		_, err := ReadJSON[BodyTest](input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		input := strings.NewReader(`string decoded as is`)
		_, err := ReadString[string](input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

type BodyTestWithInTransformer struct {
	A string
	B int
}

func (t *BodyTestWithInTransformer) InTransform() error {
	t.A = "transformed " + t.A
	return nil
}

var _ InTransformer = &BodyTestWithInTransformer{}

func TestInTransform(t *testing.T) {
	t.Run("ReadJSON", func(t *testing.T) {
		input := strings.NewReader(`{"A":"a", "B":1}`)
		body, err := ReadJSON[BodyTestWithInTransformer](input)
		require.NoError(t, err)
		require.Equal(t, BodyTestWithInTransformer{"transformed a", 1}, body)
	})
}

type transformableString string

func (t *transformableString) InTransform() error {
	*t = "transformed " + *t
	return nil
}

var _ InTransformer = new(transformableString)

func TestInTransformString(t *testing.T) {
	t.Run("ReadString", func(t *testing.T) {
		input := strings.NewReader(`coucou`)
		body, err := ReadString[transformableString](input)
		require.NoError(t, err)
		require.Equal(t, transformableString("transformed coucou"), body)
	})
}

type transformableStringWithError string

func (t *transformableStringWithError) InTransform() error {
	*t = "transformed " + *t
	return errors.New("error happened!")
}

var _ InTransformer = new(transformableStringWithError)

func TestInTransformStringWithError(t *testing.T) {
	t.Run("ReadString", func(t *testing.T) {
		input := strings.NewReader(`coucou`)
		body, err := ReadString[transformableStringWithError](input)
		require.ErrorAs(t, err, &BadRequestError{}, "Expected a BadRequestError")
		require.Equal(t, transformableStringWithError("transformed coucou"), body)
	})
}

func TestReadURLEncoded(t *testing.T) {
	t.Run("read urlencoded", func(t *testing.T) {
		input := strings.NewReader(`A=a&B=1&C=true`)
		r := httptest.NewRequest("POST", "/", input)
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		res, err := ReadURLEncoded[BodyTest](r)
		require.NoError(t, err)
		require.Equal(t, BodyTest{"a", 1, true}, res)
	})

	t.Run("read urlencoded with type error", func(t *testing.T) {
		input := strings.NewReader(`A=a&B=wrongtype&C=true`)
		r := httptest.NewRequest("POST", "/", input)
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		res, err := ReadURLEncoded[BodyTest](r)
		require.Error(t, err)
		require.Equal(t, BodyTest{"a", 0, true}, res)
	})
}
