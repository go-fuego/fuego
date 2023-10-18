package op

import (
	"errors"
	"net/http/httptest"
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
	s := NewServer()

	t.Run("can create std http handler from op controller", func(t *testing.T) {
		handler := httpHandler[ans, any](s, testController)
		if handler == nil {
			t.Error("handler is nil")
		}
	})

	t.Run("can run http handler from op controller", func(t *testing.T) {
		handler := httpHandler[ans, any](s, testController)

		req := httptest.NewRequest("GET", "/testing", nil)
		w := httptest.NewRecorder()
		handler(w, req)

		body := w.Body.String()
		require.Equal(t, crlf(`{"ans":"Hello World"}`), body)
	})

	t.Run("can handle errors in http handler from op controller", func(t *testing.T) {
		handler := httpHandler[ans, any](s, testControllerWithError)
		if handler == nil {
			t.Error("handler is nil")
		}

		req := httptest.NewRequest("GET", "/testing", nil)
		w := httptest.NewRecorder()
		handler(w, req)

		body := w.Body.String()
		require.Equal(t, crlf(`{"error":"error happened!"}`), body)
	})
}
