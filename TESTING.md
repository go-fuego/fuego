# Testing in Fuego

This guide explains how to effectively test your Fuego applications, with a focus on using the mock context for testing controllers.

## Mock Context

Fuego provides a framework-agnostic mock context that allows you to test your controllers without depending on specific web frameworks. This makes it easy to focus on testing your business logic rather than HTTP mechanics.

### Basic Usage

```go
func TestMyController(t *testing.T) {
    // Create a mock context with your request body type
    ctx := fuego.NewMockContext[UserRequest]()

    // Set up test data
    ctx.SetBody(UserRequest{
        Name: "John Doe",
        Email: "john@example.com",
    })

    // Add query parameters if needed
    ctx.SetURLValues(url.Values{
        "filter": []string{"active"},
    })

    // Add path parameters
    ctx.SetPathParam("id", "123")

    // Call your controller
    result, err := MyController(ctx)

    // Assert results
    assert.NoError(t, err)
    assert.Equal(t, expectedResult, result)
}
```

### Features

The mock context supports:

- Type-safe request bodies with generics
- URL query parameters
- Path parameters
- Headers
- Custom context values
- Request/Response objects

### Advanced Usage

#### Testing with Headers

```go
func TestControllerWithHeaders(t *testing.T) {
    ctx := fuego.NewMockContext[EmptyBody]()
    ctx.SetHeader("Authorization", "Bearer token123")
    ctx.SetHeader("Content-Type", "application/json")

    result, err := MyAuthenticatedController(ctx)
    assert.NoError(t, err)
}
```

#### Testing with Custom Context Values

```go
func TestControllerWithContext(t *testing.T) {
    ctx := fuego.NewMockContext[EmptyBody]()
    customCtx := context.WithValue(context.Background(), "user_id", "123")
    ctx.SetContext(customCtx)

    result, err := MyContextAwareController(ctx)
    assert.NoError(t, err)
}
```

#### Testing with Request/Response Objects

```go
func TestControllerWithRequestResponse(t *testing.T) {
    ctx := fuego.NewMockContext[EmptyBody]()
    w := httptest.NewRecorder()
    r := httptest.NewRequest("GET", "/test", nil)

    ctx.SetResponse(w)
    ctx.SetRequest(r)

    result, err := MyController(ctx)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, w.Code)
}
```

### Best Practices

1. **Use Table-Driven Tests**

```go
func TestUserController(t *testing.T) {
    tests := []struct {
        name    string
        body    UserRequest
        want    UserResponse
        wantErr bool
    }{
        {
            name: "valid user",
            body: UserRequest{Name: "John", Email: "john@example.com"},
            want: UserResponse{ID: "123", Name: "John"},
        },
        {
            name:    "invalid email",
            body:    UserRequest{Name: "John", Email: "invalid"},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := fuego.NewMockContext[UserRequest]()
            ctx.SetBody(tt.body)

            got, err := CreateUser(ctx)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

2. **Test Error Cases**

```go
func TestErrorHandling(t *testing.T) {
    ctx := fuego.NewMockContext[UserRequest]()
    ctx.SetBody(UserRequest{}) // Empty body should trigger validation error

    _, err := CreateUser(ctx)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "validation failed")
}
```

3. **Test Validation Rules**

```go
func TestValidation(t *testing.T) {
    ctx := fuego.NewMockContext[UserRequest]()
    ctx.SetBody(UserRequest{
        Name:  "", // Required field
        Email: "invalid-email", // Invalid format
    })

    _, err := CreateUser(ctx)
    assert.Error(t, err)
}
```

4. **Test Middleware Integration**

```go
func TestWithMiddleware(t *testing.T) {
    ctx := fuego.NewMockContext[EmptyBody]()
    ctx.SetHeader("Authorization", "Bearer valid-token")

    // Test that middleware allows the request
    result, err := AuthMiddleware(MyProtectedController)(ctx)
    assert.NoError(t, err)

    // Test that middleware blocks unauthorized requests
    ctx.SetHeader("Authorization", "invalid-token")
    _, err = AuthMiddleware(MyProtectedController)(ctx)
    assert.Error(t, err)
}
```

### Tips for Effective Testing

1. Keep tests focused on business logic
2. Use meaningful test names that describe the scenario
3. Test both success and failure cases
4. Use helper functions for common test setup
5. Test validation rules thoroughly
6. Mock external dependencies when needed
7. Use subtests for better organization
8. Test edge cases and boundary conditions

## Contributing

If you find any issues or have suggestions for improving the testing utilities, please open an issue or submit a pull request.
