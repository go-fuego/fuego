package commands

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateController(t *testing.T) {
	res, err := createNewEntityDomainFile("books", "controller.go", "booksController.go")
	require.NoError(t, err)
	require.Contains(t, res, "package books")
	require.Contains(t, res, `fuego.Get(booksGroup, "/{id}", rs.getBooks)`)
	require.Contains(t, res, `func (rs BooksResources) postBooks(c fuego.ContextWithBody[BooksCreate]) (Books, error)`)
	require.FileExists(t, "./domains/books/booksController.go")
	os.Remove("./domains/books/booksController.go")
}
