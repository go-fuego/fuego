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
	require.Equal(t, "default", tagFromType(*new(any)))
	require.Equal(t, "MyStruct", tagFromType(MyStruct{}))
	require.Equal(t, "MyStruct", tagFromType(&MyStruct{}))
	require.Equal(t, "MyStruct", tagFromType([]MyStruct{}))
	require.Equal(t, "MyStruct", tagFromType(&[]MyStruct{}))
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
