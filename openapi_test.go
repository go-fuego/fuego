package op

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type MyStruct struct {
	B string `json:"b"`
	C int    `json:"c"`
	D bool   `json:"d"`
}

func TestTagFromType(t *testing.T) {
	require.Equal(t, "unknown-interface", tagFromType(*new(any)), "behind any interface")
	require.Equal(t, "MyStruct", tagFromType(MyStruct{}))

	t.Run("behind pointers or pointers-like", func(t *testing.T) {
		require.Equal(t, "MyStruct", tagFromType(&MyStruct{}))
		require.Equal(t, "MyStruct", tagFromType([]MyStruct{}))
		require.Equal(t, "MyStruct", tagFromType(&[]MyStruct{}))
		type DeeplyNested *[]MyStruct
		require.Equal(t, "MyStruct", tagFromType(new(DeeplyNested)), "behind 4 pointers")
	})

	t.Run("safety against recursion", func(t *testing.T) {
		type DeeplyNested *[]MyStruct
		type MoreDeeplyNested *[]DeeplyNested
		require.Equal(t, "MyStruct", tagFromType(*new(MoreDeeplyNested)), "behind 5 pointers")

		require.Equal(t, "default", tagFromType(new(MoreDeeplyNested)), "behind 6 pointers")
		require.Equal(t, "default", tagFromType([]*MoreDeeplyNested{}), "behind 7 pointers")
	})
}

func TestServer_GenerateOpenAPI(t *testing.T) {
	s := NewServer()
	Get(s, "/", func(Ctx[any]) (MyStruct, error) {
		return MyStruct{}, nil
	})
	Post(s, "/post", func(Ctx[MyStruct]) ([]MyStruct, error) {
		return nil, nil
	})
	Get(s, "/post/{id}", func(Ctx[any]) (MyStruct, error) {
		return MyStruct{}, nil
	})
	require.NotPanics(t, func() {
		s.GenerateOpenAPI()
	})
}
