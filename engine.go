package fuego

func NewEngine() *Engine {
	return &Engine{
		OpenAPI:      NewOpenAPI(),
		ErrorHandler: ErrorHandler,
	}
}

// The Engine is the main struct of the framework.
type Engine struct {
	OpenAPI      *OpenAPI
	ErrorHandler func(error) error

	acceptedContentTypes []string
}
