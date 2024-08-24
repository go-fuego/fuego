package fuego

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"io"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type BodyTest struct {
	XMLName xml.Name `xml:"BodyTest"`
	A       string   `yaml:"A" xml:"A"`
	B       int      `yaml:"B" xml:"B"`
	C       bool     `yaml:"C" xml:"C"`
}

type WrongBody struct {
	XMLName xml.Name `xml:"WrongBody"`
	A       string   `yaml:"A" xml:"A"`
	B       int      `yaml:"B" xml:"B"`
}

type EOFReader struct{}

func (r *EOFReader) Read(b []byte) (int, error) {
	return 0, io.EOF
}

func TestReadJSON(t *testing.T) {
	inputStr := `{"A":"a","B":1,"C":true}`

	t.Run("ReadJSON", func(t *testing.T) {
		body, err := ReadJSON[BodyTest](context.Background(), strings.NewReader(inputStr))
		require.NoError(t, err)
		require.Equal(t, BodyTest{A: "a", B: 1, C: true}, body)
	})

	t.Run("error is EOF", func(t *testing.T) {
		body, err := ReadJSON[BodyTest](context.Background(), &EOFReader{})
		require.NoError(t, err)
		require.Equal(t, BodyTest{A: "", B: 0, C: false}, body)
	})

	t.Run("cannot read invalid JSON", func(t *testing.T) {
		_, err := ReadJSON[BodyTest](context.Background(), strings.NewReader(`{`))
		require.ErrorAs(t, err, &BadRequestError{}, "Expected a BadRequestError")
	})

	t.Run("cannot deserialize JSON to wrong struct", func(t *testing.T) {
		_, err := ReadJSON[WrongBody](context.Background(), strings.NewReader(inputStr))
		require.ErrorAs(t, err, &BadRequestError{}, "Expected a BadRequestError")
	})
}

func TestReadYAML(t *testing.T) {
	inputStr := `
A: a
B: 1
C: true
`

	t.Run("ReadYAML", func(t *testing.T) {
		body, err := ReadYAML[BodyTest](context.Background(), bytes.NewReader([]byte(inputStr)))
		require.NoError(t, err)
		require.Equal(t, BodyTest{A: "a", B: 1, C: true}, body)
	})

	t.Run("cannot read invalid YAML", func(t *testing.T) {
		_, err := ReadYAML[BodyTest](context.Background(), bytes.NewReader([]byte(`a`)))
		require.ErrorAs(t, err, &BadRequestError{}, "Expected a BadRequestError")
	})

	t.Run("cannot deserialize YAML to wrong struct", func(t *testing.T) {
		_, err := ReadYAML[WrongBody](context.Background(), bytes.NewReader([]byte(inputStr)))
		require.ErrorAs(t, err, &BadRequestError{}, "Expected a BadRequestError")
	})
}

func TestReadXML(t *testing.T) {
	inputStr := `
<BodyTest>
	<A>a</A>
	<B>1</B>
	<C>true</C>
</BodyTest>
`

	t.Run("ReadXML", func(t *testing.T) {
		body, err := ReadXML[BodyTest](context.Background(), bytes.NewReader([]byte(inputStr)))
		require.NoError(t, err)
		require.Equal(t, BodyTest{
			XMLName: xml.Name{Local: "BodyTest"},
			A:       "a",
			B:       1,
			C:       true,
		}, body)
	})

	t.Run("cannot read invalid XML", func(t *testing.T) {
		_, err := ReadXML[BodyTest](context.Background(), bytes.NewReader([]byte("<")))
		require.ErrorAs(t, err, &BadRequestError{}, "Expected a BadRequestError")
	})

	t.Run("cannot deserialize XML to wrong struct", func(t *testing.T) {
		_, err := ReadXML[WrongBody](context.Background(), bytes.NewReader([]byte(inputStr)))
		require.ErrorAs(t, err, &BadRequestError{}, "Expected a BadRequestError")
	})
}

type errorReader int

func (errorReader) Read(_ []byte) (int, error) {
	return 0, errors.New("error")
}

func TestReadString(t *testing.T) {
	t.Run("read string", func(t *testing.T) {
		input := strings.NewReader(`string decoded as is`)
		_, err := ReadString[string](context.Background(), input)
		require.NoError(t, err)
	})

	t.Run("read string alias", func(t *testing.T) {
		type StringAlias string
		input := strings.NewReader(`string decoded as is`)
		_, err := ReadString[StringAlias](context.Background(), input)
		require.NoError(t, err)
	})

	t.Run("should return error if io.ReadAll return an error", func(t *testing.T) {
		type StringAlias string
		input := errorReader(0)
		_, err := ReadString[StringAlias](context.Background(), input)
		require.Error(t, err)
	})
}

func BenchmarkReadJSON(b *testing.B) {
	for range b.N {
		input := strings.NewReader(`{"A":"a","B":1,"C":true}`)
		_, err := ReadJSON[BodyTest](context.Background(), input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadString(b *testing.B) {
	for range b.N {
		input := strings.NewReader(`string decoded as is`)
		_, err := ReadString[string](context.Background(), input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

type BodyTestWithInTransformer struct {
	A string
	B int
}

func (t *BodyTestWithInTransformer) InTransform(context.Context) error {
	t.A = "transformed " + t.A
	return nil
}

var _ InTransformer = &BodyTestWithInTransformer{}

type BodyTestWithInTransformerError struct {
	A string
	B int
}

func (t *BodyTestWithInTransformerError) InTransform(context.Context) error {
	return errors.New("error happened!")
}

var _ InTransformer = &BodyTestWithInTransformerError{}

func TestInTransform(t *testing.T) {
	t.Run("ReadJSON", func(t *testing.T) {
		input := strings.NewReader(`{"A":"a", "B":1}`)
		body, err := ReadJSON[BodyTestWithInTransformer](context.Background(), input)
		require.NoError(t, err)
		require.Equal(t, BodyTestWithInTransformer{"transformed a", 1}, body)
	})
}

type transformableString string

func (t *transformableString) InTransform(context.Context) error {
	*t = "transformed " + *t
	return nil
}

var _ InTransformer = new(transformableString)

func TestInTransformString(t *testing.T) {
	t.Run("ReadString", func(t *testing.T) {
		input := strings.NewReader(`coucou`)
		body, err := ReadString[transformableString](context.Background(), input)
		require.NoError(t, err)
		require.Equal(t, transformableString("transformed coucou"), body)
	})
}

type transformableStringWithError string

func (t *transformableStringWithError) InTransform(context.Context) error {
	*t = "transformed " + *t
	return errors.New("error happened!")
}

var _ InTransformer = new(transformableStringWithError)

func TestInTransformStringWithError(t *testing.T) {
	t.Run("ReadString", func(t *testing.T) {
		input := strings.NewReader(`coucou`)
		body, err := ReadString[transformableStringWithError](context.Background(), input)
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
		require.Equal(t, BodyTest{A: "a", B: 1, C: true}, res)
	})

	t.Run("read urlencoded with type error", func(t *testing.T) {
		input := strings.NewReader(`A=a&B=wrongtype&C=true`)
		r := httptest.NewRequest("POST", "/", input)
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		res, err := ReadURLEncoded[BodyTest](r)
		require.Error(t, err)
		require.Equal(t, BodyTest{A: "a", B: 0, C: true}, res)
	})

	t.Run("read urlencoded with transform error", func(t *testing.T) {
		input := strings.NewReader(`A=a&B=9`)
		r := httptest.NewRequest("POST", "/", input)
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		res, err := ReadURLEncoded[BodyTestWithInTransformerError](r)
		require.Error(t, err)
		require.ErrorAs(t, err, &BadRequestError{}, "Expected a BadRequestError")
		require.Equal(t, BodyTestWithInTransformerError{"a", 9}, res)
	})

	t.Run("read invalid semicolon separator in query", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/", nil)
		r.URL.RawQuery = ";invalid;"
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		_, err := ReadURLEncoded[any](r)
		require.Error(t, err)
	})
}

func TestConvertSQLNullString(t *testing.T) {
	t.Run("can convert sql.NullString", func(t *testing.T) {
		v := convertSQLNullString("test")
		require.Equal(t, "test", v.Interface().(sql.NullString).String)
	})
}

func TestConvertSQLNullBool(t *testing.T) {
	t.Run("convert sql.NullBool", func(t *testing.T) {
		v := convertSQLNullBool("false")
		require.Equal(t, false, v.Interface().(sql.NullBool).Bool)
	})

	t.Run("cannot convert sql.NullBool", func(t *testing.T) {
		v := convertSQLNullBool("hello")
		require.Equal(t, v, reflect.Value{})
	})
}
