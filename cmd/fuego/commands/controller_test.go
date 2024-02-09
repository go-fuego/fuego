package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateController(t *testing.T) {
	res, err := createController("books", "/dev/null")
	require.NoError(t, err)
	require.Contains(t, res, "package controller")
	require.Contains(t, res, `fuego.Get(booksGroup, "/{id}", rs.getBooks)`)
	require.Contains(t, res, `func (rs BooksRessources) postBooks(c *fuego.ContextWithBody[BooksCreate]) (Books, error)`)
}
