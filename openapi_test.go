package op

import (
	"fmt"
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
	Get(s, "/post/{id}", func(Ctx[any]) (MyOutputStruct, error) {
		return MyOutputStruct{}, nil
	})
	document := s.GenerateOpenAPI()
	require.NotNil(t, document)
	require.NotNil(t, document.Paths["/"])
	require.Nil(t, document.Paths.Find("/unknown"))
	require.NotNil(t, document.Paths.Find("/post"))
	require.Nil(t, document.Paths.Find("/post/{id}").Get.Responses["200"].Value.Content["application/json"].Schema.Value.Properties["unknown"])
	require.Equal(t, document.Paths.Find("/post/{id}").Get.Responses["200"].Value.Content["application/json"].Schema.Value.Properties["quantity"].Value.Type, "integer")
}

func BenchmarkRoutesRegistration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewServer(
			WithoutLogger(),
		)
		Get(s, "/", func(Ctx[any]) (MyStruct, error) {
			return MyStruct{}, nil
		})
		for j := 0; j < 100; j++ {
			Post(s, fmt.Sprintf("/post/%d", j), func(Ctx[MyStruct]) ([]MyStruct, error) {
				return nil, nil
			})
		}
		for j := 0; j < 100; j++ {
			Get(s, fmt.Sprintf("/post/{id}/%d", j), func(Ctx[any]) (MyStruct, error) {
				return MyStruct{}, nil
			})
		}
	}
}

func BenchmarkServer_GenerateOpenAPI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewServer(
			WithoutLogger(),
		)
		Get(s, "/", func(Ctx[any]) (MyStruct, error) {
			return MyStruct{}, nil
		})
		for j := 0; j < 100; j++ {
			Post(s, fmt.Sprintf("/post/%d", j), func(Ctx[MyStruct]) ([]MyStruct, error) {
				return nil, nil
			})
		}
		for j := 0; j < 100; j++ {
			Get(s, fmt.Sprintf("/post/{id}/%d", j), func(Ctx[any]) (MyStruct, error) {
				return MyStruct{}, nil
			})
		}

		s.GenerateOpenAPI()
	}
}
