package fuego

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
)

// CtxRenderer can be used with [github.com/a-h/templ]
// Example:
//
//	func getRecipes(ctx fuego.Ctx[any]) (fuego.CtxRenderer, error) {
//		recipes, err := ctx.store.GetRecipes(ctx.Context())
//		if err != nil {
//			return nil, err
//		}
//
//		return recipeComponent(recipes), nil // recipeComponent is templ component
//	}
type CtxRenderer interface {
	Render(context.Context, io.Writer) error
}

// Templ is a shortcut for [CtxRenderer], which can be used with [github.com/a-h/templ]
type Templ = CtxRenderer

// Renderer can be used with [github.com/maragudk/gomponents]
// Example:
//
//	func getRecipes(ctx fuego.Ctx[any]) (fuego.CtxRenderer, error) {
//		recipes, err := ctx.store.GetRecipes(ctx.Context())
//		if err != nil {
//			return nil, err
//		}
//
//		return recipeComponent(recipes), nil // recipeComponent is gomponents component
//	}
type Renderer interface {
	Render(io.Writer) error
}

// Gomponent is a shortcut for [Renderer], which can be used with [github.com/maragudk/gomponents]
type Gomponent = Renderer

// HTML is a marker type used to differentiate between a string response and an HTML response.
// To use templating, use [Ctx.Render].
type HTML string

// H is a shortcut for map[string]any
type H map[string]any

// loadTemplates
func (s *Server) loadTemplates(patterns ...string) error {
	tmpl, err := template.ParseFS(s.fs, patterns...)
	if err != nil {
		var pathError *fs.PathError
		if errors.As(err, &pathError) {
			wd, _ := os.Getwd()
			return fmt.Errorf("template '%s' does not exist in directory '%s': %w", pathError.Path, wd, err)
		}

		return fmt.Errorf("%w %s", err, "failed to parse templates")
	}

	s.template = tmpl

	return nil
}
