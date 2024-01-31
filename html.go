package fuego

import (
	"context"
	"fmt"
	"html/template"
	"io"
)

// CtxRenderer can be used with [github.com/a-h/templ]
// Example:
//
//	func getRecipes(ctx fuego.ContextNoBody) (fuego.CtxRenderer, error) {
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
//	func getRecipes(ctx fuego.ContextNoBody) (fuego.CtxRenderer, error) {
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
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	s.template = tmpl

	return nil
}
