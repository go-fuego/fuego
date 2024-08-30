package param_test

import (
	"testing"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func TestParam(t *testing.T) {
	t.Run("All options", func(t *testing.T) {
		s := fuego.NewServer()

		fuego.Get(s, "/test", func(c fuego.ContextNoBody) (string, error) {
			return "test", nil
		},
			option.Query("param1", "description1", param.Required(), param.Default("hey"), param.Example("example1", "you")),
			option.QueryInt("param1", "description1", param.Nullable(), param.Default(1), param.Example("example1", 1)),
			option.QueryBool("param1", "description1", param.Required(), param.Default(true), param.Example("example1", true)),
		)
	})
}
