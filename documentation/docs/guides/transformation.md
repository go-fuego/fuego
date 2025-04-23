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

	return user, nil
}
```

## Output Transformation

Output transformation is the process of transforming the data going **out of** your application before it's serialized and sent to the client.

```go
package main

import (
	"context"
	"strings"

	"github.com/go-fuego/fuego"
)

type UserResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
}

func (u *UserResponse) OutTransform(ctx context.Context) error {
	// Mask the last name for privacy
	if len(u.LastName) > 0 {
		u.LastName = u.LastName[:1] + "***"
	}

	// Generate the full name from first and last name
	u.FullName = u.FirstName + " " + u.LastName

	// Use the context to get the request ID
	if reqID, ok := RequestIDFromContext(ctx); ok {
		u.FullName += " (Request ID: " + reqID + ")"
	}

	return nil
}

var _ fuego.OutTransformer = (*UserResponse)(nil) // Ensure *UserResponse implements fuego.OutTransformer
```

In this example, the `OutTransform` method is called before the `UserResponse` struct is serialized and sent to the client. We're using it to generate a full name field and to mask the last name for privacy.

```go
func getUser(c fuego.ContextNoBody) (UserResponse, error) {
	// Get user from database or other source
	user := UserResponse{
		FirstName: "John",
		LastName:  "Doe",
	}

	// The OutTransform method will be called automatically before serialization
	// After transformation, user.FullName will be "John Doe" and user.LastName will be "D***"

	return user, nil
}
```

## Custom Validation with Transformation

Transformation methods can also be used for custom validation that goes beyond what the standard validator can do:

```go
func (u *User) InTransform(ctx context.Context) error {
	// Normalize data
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)

	// Custom validation
	if u.FirstName == "Admin" && u.LastName == "User" {
		return errors.New("reserved name combination")
	}

	return nil
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

The same principle applies to `OutTransform` for nested structs.

## Transformation Flow

Here's the complete flow of data through Fuego's transformation and validation system:

1. Request comes in with JSON/XML/etc. payload
2. Payload is deserialized into your struct
3. `InTransform` is called on your struct (if implemented)
4. Validation is performed on your struct (if validation tags are present)
5. Your controller is called with the transformed and validated struct
6. Your controller returns a response struct
7. `OutTransform` is called on your response struct (if implemented)
8. Response struct is serialized to JSON/XML/etc. and sent to the client

This flow ensures that your data is properly transformed and validated at each step of the request/response cycle.
