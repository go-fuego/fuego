package openapi3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type MyStruct struct {
	B string `json:"b"`
	C int    `json:"c"`
	D bool   `json:"d"`
}

type MyOutputStruct struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

func TestTagFromType(t *testing.T) {
	require.Equal(t, "nil", TagFromType(*new(any)), "behind any interface")
	require.Equal(t, "MyStruct", TagFromType(MyStruct{}))

	t.Run("behind pointers or pointers-like", func(t *testing.T) {
		require.Equal(t, "MyStruct", TagFromType(&MyStruct{}))
		require.Equal(t, "MyStruct", TagFromType([]MyStruct{}))
		require.Equal(t, "MyStruct", TagFromType(&[]MyStruct{}))
		type DeeplyNested *[]MyStruct
		require.Equal(t, "MyStruct", TagFromType(new(DeeplyNested)), "behind 4 pointers")
	})

	t.Run("safety against recursion", func(t *testing.T) {
		type DeeplyNested *[]MyStruct
		type MoreDeeplyNested *[]DeeplyNested
		require.Equal(t, "MyStruct", TagFromType(*new(MoreDeeplyNested)), "behind 5 pointers")

		require.Equal(t, "default", TagFromType(new(MoreDeeplyNested)), "behind 6 pointers")
		require.Equal(t, "default", TagFromType([]*MoreDeeplyNested{}), "behind 7 pointers")
	})

	t.Run("detecting string", func(t *testing.T) {
		require.Equal(t, "string", TagFromType("string"))
		require.Equal(t, "string", TagFromType(new(string)))
		require.Equal(t, "string", TagFromType([]string{}))
		require.Equal(t, "string", TagFromType(&[]string{}))
	})
}
