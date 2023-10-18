package op

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestBody struct {
	A string
	B int
	C bool
}

func TestReadJSON(t *testing.T) {
	input := strings.NewReader(`{"A":"a","B":1,"C":true}`)

	t.Run("ReadJSON", func(t *testing.T) {
		body, err := ReadJSON[TestBody](input)
		require.NoError(t, err)
		require.Equal(t, TestBody{"a", 1, true}, body)
	})

	t.Run("cannot read invalid JSON", func(t *testing.T) {
		_, err := ReadJSON[TestBody](input)
		require.Error(t, err)
	})

	t.Run("cannot deserialize JSON to wrong struct", func(t *testing.T) {
		type WrongBody struct {
			A string
			B int
			// Missing C bool
		}
		_, err := ReadJSON[WrongBody](input)
		require.Error(t, err)
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
		_, err := ReadJSON[TestBody](input)
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
