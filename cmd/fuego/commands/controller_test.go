package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateController(t *testing.T) {
	err := createController("books")
	require.NoError(t, err)
}
