# Using Fuego with Gin

Fuego can be used with Gin by using the `fuegogin` adaptor.

Instead of using the **Server** (`fuego.NewServer()`), you will use the **Engine** (`fuego.NewEngine()`) along with your router.

The usage is similar to the default server, but you will need to declare the routes with `fuegogin.Get`, `fuegogin.Post`... instead of `fuego.Get`, `fuego.Post`...

## Migrate incrementally

1. Spawn an engine with `fuego.NewEngine()`.
2. Use `fuegogin.GetGin` instead of `gin.GET` to wrap the routes with OpenAPI declaration of the route, **without even touching the existing controllers**!!!
3. Replace the controllers **one by one** with Fuego controllers. You'll get complete OpenAPI documentation, validation, Content-Negotiation for each controller you replace!
4. Enjoy the benefits of Fuego with your existing Gin application!

## Example

Please refer to the [Gin example](https://github.com/go-fuego/fuego/tree/main/examples/gin-compat) for a complete and up-to-date example.
