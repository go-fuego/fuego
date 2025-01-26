# Validation

Validation is the process of ensuring that the data provided to the application
is correct and meaningful.

With fuego, you have several options for validating data.

- struct tags with `go-validator`
- custom validation functions

## Struct tags

You can use struct tags to validate the data coming into your application.
This is a common pattern in Go and is used by many libraries,
and we use [`go-playground/validator`](https://github.com/go-playground/validator) to do so.

```go
type User struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Age       int    `json:"age" validate:"gte=0,lte=130"`
	Email     string `json:"email" validate:"email"`
}
```

## Custom validation

You can also use Fuego's [Transformation](./transformation.md) methods to validate the data.

```go
package main

import (
	"context"
	"errors"
	"strings"

	"github.com/go-fuego/fuego"
)

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (u *User) InTransform(ctx context.Context) error {
	u.FirstName = strings.ToUpper(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)

	if u.FirstName == "" {
		return errors.New("first name is required")
	}
	return nil
}

var _ fuego.InTransformer = (*User)(nil) // Ensure *User implements fuego.InTransformer
// This check is a classic example of Go's interface implementation check and we highly recommend to use it
```
