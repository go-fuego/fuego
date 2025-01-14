# Testing Fuego Controllers

Fuego provides a `MockContext` type that makes it easy to test your controllers without setting up a full HTTP server. This guide will show you how to use it effectively.

## Using MockContext

The `MockContext` type implements the `ContextWithBody` interface, allowing you to test your controllers in isolation. Here's a simple example:

```go
func TestMyController(t *testing.T) {
    // Create a new mock context with your request body type
    ctx := fuego.NewMockContext[MyRequestType]()

    // Set the request body
    ctx.SetBody(MyRequestType{
        Name: "John",
        Age:  30,
    })

    // Call your controller
    response, err := MyController(ctx)

    // Assert the results
    assert.NoError(t, err)
    assert.Equal(t, expectedResponse, response)
}
```

## Complete Example

Here's a more complete example showing how to test a controller that uses request body, query parameters, and validation:

```go
// UserSearchRequest represents the search criteria
type UserSearchRequest struct {
    MinAge    int    `json:"minAge" validate:"gte=0,lte=150"`
    MaxAge    int    `json:"maxAge" validate:"gte=0,lte=150"`
    NameQuery string `json:"nameQuery" validate:"required"`
}

// SearchUsersController is our controller to test
func SearchUsersController(c fuego.ContextWithBody[UserSearchRequest]) (UserSearchResponse, error) {
    body, err := c.Body()
    if err != nil {
        return UserSearchResponse{}, err
    }

    // Get pagination from query params
    page := c.QueryParamInt("page")
    if page < 1 {
        page = 1
    }

    // Business logic validation
    if body.MinAge > body.MaxAge {
        return UserSearchResponse{}, errors.New("minAge cannot be greater than maxAge")
    }

    // ... rest of the controller logic
}

func TestSearchUsersController(t *testing.T) {
    tests := []struct {
        name          string
        body          UserSearchRequest
        queryParams   map[string]string
        expectedError string
        expected      UserSearchResponse
    }{
        {
            name: "successful search",
            body: UserSearchRequest{
                MinAge:    20,
                MaxAge:    35,
                NameQuery: "John",
            },
            queryParams: map[string]string{
                "page": "1",
            },
            expected: UserSearchResponse{
                // ... expected response
            },
        },
        {
            name: "invalid age range",
            body: UserSearchRequest{
                MinAge:    40,
                MaxAge:    20,
                NameQuery: "John",
            },
            expectedError: "minAge cannot be greater than maxAge",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create mock context and set up the test
            ctx := fuego.NewMockContext[UserSearchRequest]()
            ctx.SetBody(tt.body)

            // Set query parameters
            if tt.queryParams != nil {
                for key, value := range tt.queryParams {
                    ctx.SetURLValues(map[string][]string{
                        key: {value},
                    })
                }
            }

            // Call the controller
            response, err := SearchUsersController(ctx)

            // Check error cases
            if tt.expectedError != "" {
                assert.EqualError(t, err, tt.expectedError)
                return
            }

            // Check success cases
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, response)
        })
    }
}
```

## Available Mock Methods

`MockContext` provides several methods to help you test different aspects of your controllers:

- `SetBody(body B)` - Set the request body
- `SetURLValues(values url.Values)` - Set query parameters
- `SetHeader(key, value string)` - Set request headers
- `SetPathParam(name, value string)` - Set path parameters
- `SetCookie(cookie http.Cookie)` - Set request cookies
- `SetContext(ctx context.Context)` - Set a custom context
- `SetResponse(w http.ResponseWriter)` - Set a custom response writer
- `SetRequest(r *http.Request)` - Set a custom request

## Best Practices

1. **Test Edge Cases**: Test both valid and invalid inputs, including validation errors.
2. **Use Table-Driven Tests**: Structure your tests as a slice of test cases for better organization.
3. **Mock Only What You Need**: Only set up the mock data that your test actually requires.
4. **Test Business Logic**: Focus on testing your business logic rather than the framework itself.
5. **Keep Tests Focused**: Each test should verify one specific behavior.

## Why Use MockContext?

We use `MockContext` with setter methods (instead of exported fields) to:

- Maintain encapsulation and consistency with real implementations
- Allow for future additions like validation or logging without breaking user code
- Ensure tests reflect how the code will behave in production

This approach makes your tests more maintainable and reliable while keeping them simple to write.
