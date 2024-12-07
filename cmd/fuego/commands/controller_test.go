package commands

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateControllerFile(t *testing.T) {
	res, err := createControllerFile("books", "controller.go", "books.go")
	require.NoError(t, err)
	require.Contains(t, res, "package controller")
	require.Contains(t, res, `fuego.Get(booksGroup, "/{id}", rs.getBooks)`)
	require.Contains(t, res, `func (rs BooksResources) postBooks(c *fuego.ContextWithBody[BooksCreate]) (Books, error)`)
	require.FileExists(t, "./controller/books.go")
	os.Remove("./controller/books.go")
}
