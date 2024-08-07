package fuego

const openapiDescription = `
This is the autogenerated OpenAPI documentation for your [Fuego](https://github.com/go-fuego/fuego) API.

Below is a Fuego Cheatsheet to help you get started. Don't hesitate to check the [Fuego documentation](https://go-fuego.github.io/fuego) for more details.

Happy coding! 🔥

## Usage

### Route registration

` + "```go" + `
func main() {
	// Create a new server
	s := fuego.NewServer()

	// Register some routes
	fuego.Post(s, "/hello", myController)
	fuego.Get(s, "/myPath", otherController)
	fuego.Put(s, "/hello", thirdController)

	adminRoutes := fuego.Group(s, "/admin")
	fuego.Use(adminRoutes, myMiddleware) // This middleware (for authentication, etc...) will be available for routes starting by /admin/*, 
	fuego.Get(adminRoutes, "/hello", groupController) // This route will be available at /admin/hello

	// Start the server
	s.Start()
}
` + "```" + `

### Basic controller

` + "```go" + `
type MyBody struct {
	Name string ` + "`json:\"name\" validate:\"required,max=30\"`" + `
}

type MyResponse struct {
	Answer string ` + "`json:\"answer\"`" + `
}

func hello(ctx *fuego.ContextWithBody[MyBody]) (*MyResponse, error) {
	body, err := ctx.Body()
	if err != nil {
		return nil, err
	}

	return &MyResponse{Answer: "Hello " + body.Name}, nil
}
` + "```" + `

### Add more details to the route

` + "```go" + `
fuego.Get(s, "/hello", myController).
	Description("This is a route that says hello").
	Summary("Say hello").
` + "```" + `
`
