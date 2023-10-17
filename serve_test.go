package op

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type ans struct {
	Ans string `json:"ans"`
}

func testController(c Ctx[any]) (ans, error) {
	return ans{Ans: "Hello World"}, nil
}

func testControllerWithError(c Ctx[any]) (ans, error) {
	return ans{}, errors.New("error happened!")
}

func TestHttpHandler(t *testing.T) {
	t.Run("can create std http handler from op controller", func(t *testing.T) {
		handler := httpHandler[ans, any](testController)
		if handler == nil {
			t.Error("handler is nil")
		}
	})

	t.Run("can run http handler from op controller", func(t *testing.T) {
		handler := httpHandler[ans, any](testController)

		req := httptest.NewRequest("GET", "/testing", nil)
		w := httptest.NewRecorder()
		handler(w, req)

		body := w.Body.String()
		require.Equal(t, body, strings.ReplaceAll(`{"ans":"Hello World"}\n`, `\n`, "\n"))
	})

	t.Run("can handle errors in http handler from op controller", func(t *testing.T) {
		handler := httpHandler[ans, any](testControllerWithError)
		if handler == nil {
			t.Error("handler is nil")
		}

		req := httptest.NewRequest("GET", "/testing", nil)
		w := httptest.NewRecorder()
		handler(w, req)

		body := w.Body.String()
		require.Equal(t, body, strings.ReplaceAll(`{"error":"error happened!"}\n`, `\n`, "\n"))
	})
}
