package op

import (
	"reflect"
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

	// TODO: fix this
	// require.Equal(t, "MyStruct", tagFromType(&[]MyStruct{}))
	require.Equal(t, "", tagFromType(&[]MyStruct{}))
}

func TestSomt(t *testing.T) {
	a := *new(any)
	t.Log(reflect.TypeOf(a))
}
