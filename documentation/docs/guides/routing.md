# Routing

## Basic Routing

Fuego uses the standard HTTP methods to define routes: `fuego.Get`, `fuego.Post`, `fuego.Put`, `fuego.Patch`, `fuego.Delete`, `fuego.Options`, `fuego.Head` and `fuego.All`.

```go
package main

import (
	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/books", listBooks)
	fuego.Post(s, "/books", createBook)
	fuego.Get(s, "/books/{id}", getBook)
	fuego.Put(s, "/books/{id}", updateBook)
	fuego.Delete(s, "/books/{id}", deleteBook)

	s.Run()
}
```

## Route Groups

You can group routes with a common prefix using the `fuego.Group` function. This is useful for organizing your routes and applying common middleware or options to a group of routes.

```go
package main

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func main() {
	s := fuego.NewServer()

	// Create a group for API routes
	api := fuego.Group(s, "/api",
		option.Tags("API"), // Apply tags to all routes in this group
	)

	// Create a subgroup for book-related routes
	books := fuego.Group(api, "/books",
		option.Tags("Books"), // Apply additional tags to this subgroup
	)

	// Routes will be available at /api/books
	fuego.Get(books, "/", listBooks)
	fuego.Post(books, "/", createBook)
	fuego.Get(books, "/{id}", getBook)
	fuego.Put(books, "/{id}", updateBook)
	fuego.Delete(books, "/{id}", deleteBook)

	// Routes will be available at /api/users
	users := fuego.Group(api, "/users")
	fuego.Get(users, "/", listUsers)
	fuego.Get(users, "/{id}", getUser)

	s.Run()
}
```

## Path Parameters

Fuego supports path parameters using the `{paramName}` syntax. You can access these parameters in your controller using the `c.PathParam` method.

```go
package main

import (
	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/books/{id}", getBook)
	fuego.Get(s, "/users/{userId}/books/{bookId}", getUserBook)

	s.Run()
}

func getBook(c fuego.ContextNoBody) (string, error) {
	id := c.PathParam("id")
	return "Book ID: " + id, nil
}

func getUserBook(c fuego.ContextNoBody) (string, error) {
	userId := c.PathParam("userId")
	bookId := c.PathParam("bookId")
	return "User ID: " + userId + ", Book ID: " + bookId, nil
}
```

## Wildcard Parameters

You can use the `{param...}` syntax to match all remaining segments of a path. This is useful for creating catch-all routes.

```go
package main

import (
	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/files/{path...}", getFile)

	s.Run()
}

func getFile(c fuego.ContextNoBody) (string, error) {
	path := c.PathParam("path")
	return "File path: " + path, nil
}
```

## Standard HTTP Handlers

If you prefer to use standard `http.Handler` functions, Fuego provides `GetStd`, `PostStd`, etc. methods that accept standard HTTP handlers.

```go
package main

import (
	"net/http"

	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()

	fuego.GetStd(s, "/standard", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from standard HTTP handler"))
	})

	s.Run()
}
```

## Route Options

You can add options to routes to provide additional information for OpenAPI documentation or to add middleware. See the [Options](./options.md) guide for more details.

```go
package main

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/books", listBooks,
		option.Summary("List all books"),
		option.Description("Returns a list of all books in the library"),
		option.Tags("Books"),
		option.QueryInt("page", "Page number", param.Default(1)),
		option.QueryInt("limit", "Items per page", param.Default(10)),
	)

	s.Run()
}
```
