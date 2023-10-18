package op

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// crlf adds a crlf to the end of a string.
func crlf(s string) string {
	return s + "\n"
}

type response struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func TestJSON(t *testing.T) {
	t.Run("can serialize json", func(t *testing.T) {
		w := httptest.NewRecorder()
		SendJSON(w, response{Message: "Hello World", Code: 200})
		body := w.Body.String()

		require.Equal(t, crlf(`{"message":"Hello World","code":200}`), body)
	})
}

func TestXML(t *testing.T) {
	t.Run("can serialize xml", func(t *testing.T) {
		w := httptest.NewRecorder()
		SendXML(w, response{Message: "Hello World", Code: 200})
		body := w.Body.String()

		require.Equal(t, `<response><Message>Hello World</Message><Code>200</Code></response>`, body)
	})

	t.Run("can serialize xml error", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := ErrorResponse{Message: "Hello World"}
		SendXMLError(w, err)
		body := w.Body.String()

		require.Equal(t, `<ErrorResponse><Error>Hello World</Error></ErrorResponse>`, body)
	})
}
