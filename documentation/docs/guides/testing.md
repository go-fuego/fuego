# Testing Fuego Controllers

Fuego provides a `MockContext` type that makes it easy to test your controllers without setting up a full HTTP server. This guide will show you how to use it effectively.

## Using MockContext

The `MockContext` type implements the `ContextWithBody` interface, allowing you to test your controllers in isolation. Here's a simple example:

```go
func TestMyController(t *testing.T) {
    // Create a new mock context with your request body type
    ctx := fuego.NewMockContext[MyRequestType]()

    // Set the request body
    ctx.BodyData = MyRequestType{
        Name: "John",
        Age:  30,
    }

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
            ctx.BodyData = tt.body

            // Set up OpenAPI parameters for validation
            ctx.OpenAPIParams = map[string]internal.OpenAPIParam{
                "page": {
                    Name:        "page",
                    Description: "Page number",
                    Type:        fuego.QueryParamType,
                    GoType:      "integer",
                    Default:     1,
                },
                // Add other OpenAPI parameters as needed
            }

            // Set query parameters
            if tt.queryParams != nil {
                values := make(url.Values)
                for key, value := range tt.queryParams {
                    values.Set(key, value)
                }
                ctx.UrlValues = values
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

## Available Fields and Methods

`MockContext` provides public fields for easy testing:

- `BodyData` - The request body
- `HeadersData` - HTTP headers
- `PathParamsData` - Path parameters
- `ResponseData` - Response writer
- `RequestData` - HTTP request
- `CookiesData` - HTTP cookies
- `UrlValues` (from CommonContext) - Query parameters
- `OpenAPIParams` (from CommonContext) - OpenAPI parameter definitions

And implements all the methods from `ContextWithBody`:

- `Body()` - Get the request body
- `MustBody()` - Get the request body (panics on error)
- `Header(key string)` - Get request header
- `HasHeader(key string)` - Check if header exists
- `PathParam(name string)` - Get path parameter
- `Cookie(name string)` - Get request cookie
- `HasCookie(key string)` - Check if cookie exists
- `Response()` - Get response writer
- `Request()` - Get request

Additionally, since `MockContext` embeds `CommonContext`, you get access to all the common functionality:

- `QueryParam(name string)` - Get query parameter
- `QueryParamInt(name string)` - Get query parameter as int
- `QueryParamBool(name string)` - Get query parameter as bool
- `QueryParamArr(name string)` - Get query parameter as string array
- `HasQueryParam(name string)` - Check if query parameter exists

## OpenAPI Parameter Setup

When using query parameters, you should define their OpenAPI specifications to avoid warnings and ensure proper validation:

```go
ctx.OpenAPIParams = map[string]internal.OpenAPIParam{
    "page": {
        Name:        "page",
        Description: "Page number",
        Type:        fuego.QueryParamType,
        GoType:      "integer",
        Default:     1,
    },
    "perPage": {
        Name:        "perPage",
        Description: "Items per page",
        Type:        fuego.QueryParamType,
        GoType:      "integer",
        Default:     20,
    },
}
```

Available OpenAPI parameter fields:

- `Name` - Parameter name
- `Description` - Parameter description
- `Type` - Parameter type (QueryParamType, HeaderParamType, CookieParamType)
- `GoType` - Go type ("integer", "string", "boolean")
- `Default` - Default value
- `Required` - Whether the parameter is required
- `Nullable` - Whether the parameter can be null
- `Examples` - Example values for documentation

## Best Practices

1. **Test Edge Cases**: Test both valid and invalid inputs, including validation errors.
2. **Use Table-Driven Tests**: Structure your tests as a slice of test cases for better organization.
3. **Mock Only What You Need**: Only set up the mock data that your test actually requires.
4. **Test Business Logic**: Focus on testing your business logic rather than the framework itself.
5. **Keep Tests Focused**: Each test should verify one specific behavior.
6. **Define OpenAPI Parameters**: Always define OpenAPI parameters for query parameters to ensure proper validation.

## Why Use MockContext?

The `MockContext` implementation:

- Embeds `CommonContext` to provide consistent behavior with real implementations
- Provides public fields for easy testing while maintaining interface compatibility
- Handles OpenAPI validation and parameter processing through CommonContext
- Makes tests easy to write and maintain

This approach makes your tests more maintainable and reliable while keeping them simple to write.
