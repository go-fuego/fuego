# CRUD

How to write simple CRUD operations.

This tutorial relies on the previous Hello world tutorial [Start from scratch.](/docs/tutorials/hello-world#start-from-scratch)

## Generation

```bash
go install github.com/go-fuego/fuego/cmd/fuego@latest
fuego controller books

# or in one line:
go run github.com/go-fuego/fuego/cmd/fuego@latest controller books
```

After executing the above code,
you need to slightly modify the generated controllers/books.go and main.go files.

```go title="controllers/books.go" {8-9,28-39}
package controller

import (
	"github.com/go-fuego/fuego"
)

type BooksRessources struct {
	// CHANGE NEXT LINE (BooksService -> RealBooksService)
	BooksService RealBooksService
}

type Books struct {
	ID   string `json:"id"`


// ....
// ....

type BooksService interface {
	GetBooks(id string) (Books, error)
	CreateBooks(BooksCreate) (Books, error)
	GetAllBooks() ([]Books, error)
	UpdateBooks(id string, input BooksUpdate) (Books, error)
	DeleteBooks(id string) (any, error)
}


// ADD THIS CODE BLOCK
type RealBooksService struct {
	BooksService
}

func (s RealBooksService) GetBooks(id string) (Books, error) {
	return Books{
		ID:   id,
		Name: "Test book data",
	}, nil
}

// TODO: Other BooksService interface implementations

// END OF CODE BLOCK
```

```go title="main.go" {4-5,13-14}
package main

import (
	// ADD NEXT LINE
	"hello-fuego/controllers"
	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()
	// ....

	// ADD NEXT LINE
	controllers.BooksRessources{}.Routes(s)
	s.Run()
}
```

If you've followed this far, `/books/:id` (GetBooks) has been implemented. 🥳

The generator will create the following routes:

- `GET /books`: list all books
- `POST /books`: create a new book
- `GET /books/:id`: get a book by id
- `PUT /books/:id`: update a book by id
- `PATCH /books/:id`: update a book by id
- `DELETE /books/:id`: delete a book by id

## Manual

:::tip

Fuego comes with a [**generator**](#generation)
that can generates CRUD routes and controllers for you!

:::

```go title="main.go"
package main

import (
	"hello-fuego/controllers"
	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()

	// List all books
	fuego.Get(s, "/books", controllers.GetBooks)

	// Create a new book
	fuego.Post(s, "/books", controllers.CreateBook)

	// Get a book by id
	fuego.Get(s, "/books/:id", controllers.GetBook)

	// Update a book by id
	fuego.Put(s, "/books/:id", controllers.UpdateBook)

	// Update a book by id
	fuego.Patch(s, "/books/:id", controllers.UpdateBook)

	// Delete a book by id
	fuego.Delete(s, "/books/:id", controllers.DeleteBook)

	s.Run()
}
```

```go title="controllers/books.go"
package controllers

import (
	"github.com/go-fuego/fuego"
)

type Book struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type BookToCreate struct {
	Title string `json:"title"`
}

func GetBooks(c *fuego.ContextNoBody) ([]Book, error) {
	// Your code here
	return nil, nil
}

func CreateBook(c *fuego.ContextWithBody[BookToCreate]) (Book, error) {
	// Your code here
	return Book{}, nil
}

func GetBook(c *fuego.ContextNoBody) (Book, error) {
	// Your code here
	return Book{}, nil
}

func UpdateBook(c *fuego.ContextWithBody[Book]) (Book, error) {
	// Your code here
	return Book{}, nil
}

func DeleteBook(c *fuego.ContextNoBody) error {
	// Your code here
	return nil
}

```
