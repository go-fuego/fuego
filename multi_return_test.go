package fuego

import (
	"context"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockCtxRenderer struct {
	RenderFunc func(context.Context, io.Writer) error
}

var _ CtxRenderer = MockCtxRenderer{}

func (m MockCtxRenderer) Render(c context.Context, w io.Writer) error {
	if m.RenderFunc == nil {
		return errors.New("RenderFunc is nil")
	}
	return m.RenderFunc(c, w)
}

type MyType struct {
	Name string
}

func TestMultiReturn(t *testing.T) {
	s := NewServer()

	Get(s, "/data", func(c ContextNoBody) (DataOrTemplate[MyType], error) {
		entity := MyType{Name: "Ewen"}

		return DataOrTemplate[MyType]{
			Data: entity,
			Template: MockCtxRenderer{
				RenderFunc: func(c context.Context, w io.Writer) error {
					_, err := w.Write([]byte(`<div>` + entity.Name + `</div>`))
					return err
				},
			},
		}, nil
	})

	Get(s, "/other", func(c ContextNoBody) (DataOrTemplate[MyType], error) {
		entity := MyType{Name: "Ewen"}

		return DataOrHTML(
			entity,
			MockCtxRenderer{
				RenderFunc: func(c context.Context, w io.Writer) error {
					_, err := w.Write([]byte(`<div>` + entity.Name + `</div>`))
					return err
				},
			}), nil
	})

	t.Run("requests HTML by default", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/data", nil)

		s.Mux.ServeHTTP(recorder, req)

		require.Equal(t, 200, recorder.Code)
		require.Equal(t, "text/html; charset=utf-8", recorder.Header().Get("Content-Type"))
		require.Equal(t, `<div>Ewen</div>`, recorder.Body.String())
	})

	t.Run("requests JSON", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/data", nil)
		req.Header.Set("Accept", "application/json")

		s.Mux.ServeHTTP(recorder, req)

		require.Equal(t, 200, recorder.Code)
		require.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		require.Equal(t, crlf(`{"Name":"Ewen"}`), recorder.Body.String())
	})

	t.Run("requests XML", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/data", nil)
		req.Header.Set("Accept", "application/xml")

		s.Mux.ServeHTTP(recorder, req)

		require.Equal(t, 200, recorder.Code)
		require.Equal(t, "application/xml", recorder.Header().Get("Content-Type"))
		require.Equal(t, `<MyType><Name>Ewen</Name></MyType>`, recorder.Body.String())
	})

	t.Run("requests HTML", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/data", nil)
		req.Header.Set("Accept", "text/html")

		s.Mux.ServeHTTP(recorder, req)

		require.Equal(t, 200, recorder.Code)
		require.Contains(t, recorder.Header().Get("Content-Type"), "text/html")
		require.Equal(t, `<div>Ewen</div>`, recorder.Body.String())
	})
}
