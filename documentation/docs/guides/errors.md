# Error handling

Error handling is a crucial part of any application. It is important to handle errors gracefully and provide meaningful feedback to the user. In this guide, we will cover how to handle errors in a Fuego application.

## Error handling in Fuego

Fuego [controllers](./controllers) returns a value and an error. If the error is not `nil`, it means that an error occurred while processing the request. The error will be returned to the client as a JSON response.

By default, Fuego implements [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457), which defines a standard error format for HTTP APIs. We strongly recommend following this standard, but you can also use your own errors.

The error type returned as JSON is `fuego.HTTPError`. It has a `Status` and a `Info` field. The `Status` field is an integer that represents the error code. The `Info` field is a string that contains a human-readable error message.

If your error implements `Status() int` and `Info()` methods, the error will include the status code and the error message in the `fuego.HTTPError` response.

```go
type MyCustomError struct {
	Status int
	Message string
}

func (e MyCustomError) Status() int {
	return e.Status
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
