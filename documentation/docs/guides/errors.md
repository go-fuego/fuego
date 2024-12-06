# Error handling

Error handling is a crucial part of any application. It is important to handle errors gracefully and provide meaningful feedback to the user. In this guide, we will cover how to handle errors in a Fuego application.

## Error handling in Fuego

Fuego [controllers](./controllers) returns a value and an error. If the error is not `nil`,
it means that an error occurred while processing the request.
The error will be returned to the client as a JSON or XML response.

The default error handler will transform any error that implements the
`fuego.ErrorWithStatus` or `fuego.ErrorWithDetail` interfaces into a `fuego.HTTPError`. The `fuego.HTTPError` implements
[RFC 9457](https://www.rfc-editor.org/rfc/rfc9457), which defines a standard error format for HTTP APIs.
We strongly recommend following this standard, but you can also use your own errors.

The default `fuego.ErrorHandler` can be overridden using `fuego.WithErrorHandler` at fuego Server creation time.

The error type of `fuego.HTTPError` is returned as JSON or XML depending on the content-type specified.
It's structure is the following:

```go
// HTTPError is the error response used by the serialization part of the framework.
type HTTPError struct {
	// Developer readable error message. Not shown to the user to avoid security leaks.
	Err error `json:"-" xml:"-"`
	// URL of the error type. Can be used to lookup the error in a documentation
	Type string `json:"type,omitempty" xml:"type,omitempty" description:"URL of the error type. Can be used to lookup the error in a documentation"`
	// Short title of the error
	Title string `json:"title,omitempty" xml:"title,omitempty" description:"Short title of the error"`
	// HTTP status code. If using a different type than [HTTPError], for example [BadRequestError],
	// this will be automatically overridden after Fuego error handling.
	Status int `json:"status,omitempty" xml:"status,omitempty" description:"HTTP status code" example:"403"`
	// Human readable error message
	Detail   string      `json:"detail,omitempty" xml:"detail,omitempty" description:"Human readable error message"`
	Instance string      `json:"instance,omitempty" xml:"instance,omitempty"`
	Errors   []ErrorItem `json:"errors,omitempty" xml:"errors,omitempty"`
}
```

If your error implements `fuego.ErrorWithStatus` or `fuego.ErrorWithDetail`
the error will be returned as a `fuego.HTTPError`.

```go
// ErrorWithStatus is an interface that can be implemented by an error to provide
// a status code
type ErrorWithStatus interface {
	error
	StatusCode() int
}

// ErrorWithDetail is an interface that can be implemented by an error to provide
// an additional detail message about the error
type ErrorWithDetail interface {
	error
	DetailMsg() string
}
```

Example:

```go
type MyCustomError struct {
	Err     error  `json:"error"`
	Message string `json:"message"`
}

var _ fuego.ErrorWithStatus = MyCustomError{}
var _ fuego.ErrorWithDetail = MyCustomError{}

func (e MyCustomError) Error() string { return e.Err.Error() }

func (e MyCustomError) StatusCode() int { return http.StatusTeapot }

func (e MyCustomError) DetailMsg() string {
	return strings.Split(e.Error(), " ")[1]
}
```

Alternatively, you can always use `fuego.HTTPError` directly such as:

```go
err := fuego.HTTPError{
	Title:  "unauthorized access",
	Detail: "wrong username or password",
	Status: http.StatusUnauthorized,
}
```

## Default errors

Fuego provides a set of default errors that you can use in your application.

- `fuego.BadRequestError`: 400 Bad Request
- `fuego.UnauthorizedError`: 401 Unauthorized
- `fuego.ForbiddenError`: 403 Forbidden
- `fuego.NotFoundError`: 404 Not Found
- `fuego.ConflictError`: 409 Conflict
- `fuego.InternalServerError`: 500 Internal Server Error
- `fuego.NotAcceptableError`: 406 Not Acceptable
