# CRUD

How to write simple CRUD operations.

:::tip

Fuego comes with a **generator** that can generates CRUD routes and controllers for you!

:::

## Generation

```bash
go install github.com/go-fuego/fuego/cmd/fuego@latest
fuego controller books

# or in one line:
go run github.com/go-fuego/fuego/cmd/fuego@latest controller books
```

## The routes

The generator will create the following routes:

- `GET /books`: list all books
- `POST /books`: create a new book
- `GET /books/:id`: get a book by id
- `PUT /books/:id`: update a book by id
- `PATCH /books/:id`: update a book by id
- `DELETE /books/:id`: delete a book by id

If you want to create the routes manually, you can use the following code:

```go title="main.go"
package main

import (
	"github.com/go-fuego/fuego"
)

type Book struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func main() {
	s := fuego.NewServer()

	// List all books
	fuego.Get(s, "/books", getBooks)

	// Create a new book
	fuego.Post(s, "/books", createBook)

	// Get a book by id
	fuego.Get(s, "/books/:id", getBook)

	// Update a book by id
	fuego.Put(s, "/books/:id", updateBook)

	// Update a book by id
	fuego.Patch(s, "/books/:id", updateBook)

	// Delete a book by id
	fuego.Delete(s, "/books/:id", deleteBook)

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

func getBooks(c fuego.ContextNoBody) ([]Book, error) {
	// Your code here
	return nil, nil
}

func createBook(c fuego.ContextWithBody[BookToCreate]) (Book, error) {
	// Your code here
	return Book{}, nil
}

func getBook(c fuego.ContextNoBody) (Book, error) {
	// Your code here
	return Book{}, nil
}

func updateBook(c fuego.ContextWithBody[Book]) (Book, error) {
	// Your code here
	return Book{}, nil
}

func deleteBook(c fuego.ContextNoBody) error {
	// Your code here
	return nil
}

```
