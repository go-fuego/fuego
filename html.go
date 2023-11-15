package op

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

// LoadTemplates
// Deprecated: Just not used.
func (s *Server) LoadTemplates(myFs fs.FS, patterns ...string) error {
	if myFs == nil {
		myFs = os.DirFS("./templates")
	}
	s.fs = myFs
	if len(patterns) == 0 {
		patterns = []string{"**/*.html"}
	}
	tmpl, err := template.ParseFS(myFs, patterns...)
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
