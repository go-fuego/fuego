package fuego

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
)

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
