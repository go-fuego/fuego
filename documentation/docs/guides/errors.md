# Error handling

Error handling is a crucial part of any application. It is important to handle errors gracefully and provide meaningful feedback to the user. In this guide, we will cover how to handle errors in a Fuego application.

## Default error handling in Fuego

Fuego [controllers](./controllers) returns a value and an error. If the error is not `nil`,
it means that an error occurred while processing the request.
The error will be returned to the client as a JSON or XML response, and logged to the console.

```go
func MyController(c fuego.ContextNoBody) (string, error) {
	return "", errors.New("an error occurred") // Sends a 500. Will be logged but not returned to the client as it is not serializable.
}
```

You are encouraged to use provided fuego error types, that implement [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457), which defines a standard error format for HTTP APIs.

```go
func MyController(c fuego.ContextNoBody) (string, error) {
	_, err := someFunction()
	if err != nil {
		return "", fuego.BadRequestError{Title: "You cannot do that", Err: err} // Returns and logs a structured 400 error.
	}

	return "success", nil
}
```

- `fuego.BadRequestError`: 400 Bad Request
- `fuego.UnauthorizedError`: 401 Unauthorized
- `fuego.ForbiddenError`: 403 Forbidden
- `fuego.NotFoundError`: 404 Not Found
- `fuego.NotAcceptableError`: 406 Not Acceptable
- `fuego.ConflictError`: 409 Conflict
- `fuego.InternalServerError`: 500 Internal Server Error

## Custom error types

The default error handler will transform any error that implements the
`fuego.ErrorWithStatus` interfaces into a `fuego.HTTPError` (that implements
[RFC 9457](https://www.rfc-editor.org/rfc/rfc9457)).

The error type of `fuego.HTTPError` is returned as JSON or XML depending on the `Accept` header specified.
It's structure is the following:

```go
type ErrorWithStatus interface {
	error
	StatusCode() int
}
```

Example:

```go
type MyCustomError struct {
	Err     error  `json:"error"`
	Message string `json:"message"`
}

var (
	_ fuego.ErrorWithStatus = MyCustomError{}
	_ fuego.ErrorWithDetail = MyCustomError{}
)

func (e MyCustomError) Error() string { return e.Err.Error() }

func (e MyCustomError) StatusCode() int { return http.StatusTeapot }

func (e MyCustomError) DetailMsg() string {
	return strings.Split(e.Error(), " ")[1]
}
```

Alternatively, you can always use `fuego.HTTPError` directly such as:

```go
err := fuego.HTTPError{
	Title:  "Custom error",
	Detail: "This is a custom error",
	Status: http.StatusTeapot,
}
```

## Custom error handling

The default `fuego.ErrorHandler` can be overridden using `fuego.WithErrorHandler` at fuego `Engine` creation time. Example mapping sqlite errors to HTTP errors.

```go
import (
	"errors"

	"modernc.org/sqlite"

	"github.com/go-fuego/fuego"
)

func sqliteErrorHandler(err error) error {
	var sqliteError *sqlite.Error
	if errors.As(err, &sqliteError) {
		sqliteErrorCode := sqliteError.Code()
		switch sqliteErrorCode {
		case 1555, 2067 /* UNIQUE constraint failed */ :
			return fuego.ConflictError{Title: "Duplicate", Detail: sqliteError.Error(), Err: sqliteError}
		default:
			return fuego.InternalServerError{Title: "Internal Server Error", Detail: sqliteError.Error(), Err: sqliteError}
		}
	}

	return err
}

// apply all error handlers that we want, like a chain of middlewares
func customErrorHandler(err error) error {
	return fuego.ErrorHandler(sqliteErrorHandler(myOtherCustomErrorHandler(err)))
}

func main() {
	s := fuego.NewServer(
		fuego.WithEngineOptions(
			fuego.WithErrorHandler(customErrorHandler))
		),
	)

	fuego.Get(s, "/", MyController)

	s.Run()
}
```
