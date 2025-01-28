# Transformation

With Fuego, you can transform data coming in and out of your application. This is useful for a variety of reasons, such as:

- Converting data from one format to another
- Normalizing data
- Sanitizing data
- Masking sensitive data
- And more...

## Input Transformation

Input transformation is the process of transforming the data coming **into** your application.

```go
package main

import (
	"context"
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
	return nil
}

var _ fuego.InTransformer = (*User)(nil) // Ensure *User implements fuego.InTransformer
// This check is a classic example of Go's interface implementation check and we highly recommend to use it
```

In the example above, we have a `User` struct with two fields: `FirstName` and `LastName`. We also have a method called `InTransform` that takes a `context.Context` and returns an `error`. This method is called before the data is unmarshaled into the `User` struct. In this method, we are transforming the `FirstName` to uppercase and trimming any whitespace from the `LastName`.

When using the following controller, the data will be transformed before it is unmarshaled into the `User` struct.

```go
func echoCapitalized(c fuego.ContextWithBody[User]) (User, error) {
	user, err := c.Body()
	if err != nil {
		return User{}, err
	}

	// user.FirstName is in uppercase

	return u, nil
}
```

## Recursion

Transformation is **not recursive**. If you have nested structs, you will need to transform each struct individually. This is done **on purpose** to give you more control over the transformation process: no assumptions are made about how you want to transform your data, no "magic".

```go
type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
}

func (a *Address) InTransform(ctx context.Context) error {
	a.Street = strings.TrimSpace(a.Street)
	a.City = strings.ToUpper(a.City)
	return nil
}

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	// highlight-next-line
	Address Address `json:"address"` // Nested struct
}

func (u *User) InTransform(ctx context.Context) error {
	u.FirstName = strings.ToUpper(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)

	// highlight-next-line
	err := u.Address.InTransform(ctx) // Transform the nested struct
	if err != nil {
		return err
	}

	return nil
}
```
