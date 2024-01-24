package fuego

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// dummyMiddleware sets the X-Test header on the request and the X-Test-Response header on the response.
func dummyMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-Test", "test")
		w.Header().Set("X-Test-Response", "response")
		handler.ServeHTTP(w, r)
	})
}

func TestUseStd(t *testing.T) {
	s := NewServer()
	UseStd(s, dummyMiddleware)
	GetStd(s, "/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test") != "test" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("middleware not registered"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test successful"))
	})

	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	w := httptest.NewRecorder()

	s.Mux.ServeHTTP(w, r)

	require.Equal(t, w.Code, http.StatusOK)
	require.Equal(t, w.Body.String(), "test successful")
}

func TestAll(t *testing.T) {
	s := NewServer()
	All(s, "/test", func(ctx *ContextNoBody) (string, error) {
		return "test", nil
	})

	t.Run("get", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, w.Code, http.StatusOK)
		require.Equal(t, w.Body.String(), "\"test\"\n")
	})

	t.Run("post", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/test", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, w.Code, http.StatusOK)
		require.Equal(t, w.Body.String(), "\"test\"\n")
	})
}

func TestGet(t *testing.T) {
	s := NewServer()
	Get(s, "/test", func(ctx *ContextNoBody) (string, error) {
		return "test", nil
	})

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	s.Mux.ServeHTTP(w, r)

	require.Equal(t, w.Code, http.StatusOK)
	require.Equal(t, w.Body.String(), "\"test\"\n")
}

func TestPost(t *testing.T) {
	s := NewServer()
	Post(s, "/test", func(ctx *ContextNoBody) (string, error) {
		return "test", nil
	})

	r := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()

	s.Mux.ServeHTTP(w, r)

	require.Equal(t, w.Code, http.StatusOK)
	require.Equal(t, w.Body.String(), "\"test\"\n")
}

func TestPut(t *testing.T) {
	s := NewServer()
	Put(s, "/test", func(ctx *ContextNoBody) (string, error) {
		return "test", nil
	})

	r := httptest.NewRequest(http.MethodPut, "/test", nil)
	w := httptest.NewRecorder()

	s.Mux.ServeHTTP(w, r)

	require.Equal(t, w.Code, http.StatusOK)
	require.Equal(t, w.Body.String(), "\"test\"\n")
}

func TestPatch(t *testing.T) {
	s := NewServer()
	Patch(s, "/test", func(ctx *ContextNoBody) (string, error) {
		return "test", nil
	})

	r := httptest.NewRequest(http.MethodPatch, "/test", nil)
	w := httptest.NewRecorder()

	s.Mux.ServeHTTP(w, r)

	require.Equal(t, w.Code, http.StatusOK)
	require.Equal(t, w.Body.String(), "\"test\"\n")
}

func TestDelete(t *testing.T) {
	s := NewServer()
	Delete(s, "/test", func(ctx *ContextNoBody) (string, error) {
		return "test", nil
	})

	r := httptest.NewRequest(http.MethodDelete, "/test", nil)
	w := httptest.NewRecorder()

	s.Mux.ServeHTTP(w, r)

	require.Equal(t, w.Code, http.StatusOK)
	require.Equal(t, w.Body.String(), "\"test\"\n")
}

func TestHandle(t *testing.T) {
	s := NewServer()
	Handle(s, "/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test successful"))
	}))

	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	w := httptest.NewRecorder()

	s.Mux.ServeHTTP(w, r)

	require.Equal(t, w.Code, http.StatusOK)
	require.Equal(t, w.Body.String(), "test successful")
}

func TestGetStd(t *testing.T) {
	s := NewServer()
	GetStd(s, "/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test successful"))
	})

	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	w := httptest.NewRecorder()

	s.Mux.ServeHTTP(w, r)

	require.Equal(t, w.Code, http.StatusOK)
	require.Equal(t, w.Body.String(), "test successful")
}

func TestPostStd(t *testing.T) {
	s := NewServer()
	PostStd(s, "/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test successful"))
	})

	r := httptest.NewRequest(http.MethodPost, "/test", nil)

	w := httptest.NewRecorder()

	s.Mux.ServeHTTP(w, r)

	require.Equal(t, w.Code, http.StatusOK)
	require.Equal(t, w.Body.String(), "test successful")
}

func TestPutStd(t *testing.T) {
	s := NewServer()
	PutStd(s, "/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test successful"))
	})

	r := httptest.NewRequest(http.MethodPut, "/test", nil)

	w := httptest.NewRecorder()

	s.Mux.ServeHTTP(w, r)

	require.Equal(t, w.Code, http.StatusOK)
	require.Equal(t, w.Body.String(), "test successful")
}

func TestPatchStd(t *testing.T) {
	s := NewServer()
	PatchStd(s, "/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test successful"))
	})

	r := httptest.NewRequest(http.MethodPatch, "/test", nil)

	w := httptest.NewRecorder()

	s.Mux.ServeHTTP(w, r)

	require.Equal(t, w.Code, http.StatusOK)
	require.Equal(t, w.Body.String(), "test successful")
}

func TestDeleteStd(t *testing.T) {
	s := NewServer()
	DeleteStd(s, "/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test successful"))
	})

	r := httptest.NewRequest(http.MethodDelete, "/test", nil)

	w := httptest.NewRecorder()

	s.Mux.ServeHTTP(w, r)

	require.Equal(t, w.Code, http.StatusOK)
	require.Equal(t, w.Body.String(), "test successful")
}

func TestSetTags(t *testing.T) {
	s := NewServer()
	route := Get(s, "/test", func(ctx *ContextNoBody) (string, error) {
		return "test", nil
	}).
		SetTags("my-tag").
		WithDescription("my description").
		WithSummary("my summary").
		SetDeprecated()

	require.Equal(t, route.operation.Tags, []string{"my-tag"})
	require.Equal(t, route.operation.Description, "my description")
	require.Equal(t, route.operation.Summary, "my summary")
	require.Equal(t, route.operation.Deprecated, true)
}

func TestAddTags(t *testing.T) {
	s := NewServer()
	route := Get(s, "/test", func(ctx *ContextNoBody) (string, error) {
		return "test", nil
	}).
		AddTags("my-tag").
		AddTags("my-other-tag")

	require.Equal(t, route.operation.Tags, []string{"string", "my-tag", "my-other-tag"})
}

func TestRemoveTags(t *testing.T) {
	s := NewServer()
	route := Get(s, "/test", func(ctx *ContextNoBody) (string, error) {
		return "test", nil
	}).
		AddTags("my-tag").
		RemoveTags("my-tag", "string").
		AddTags("my-other-tag")

	require.Equal(t, route.operation.Tags, []string{"my-other-tag"})
}

func TestWithQueryParams(t *testing.T) {
	s := NewServer()
	route := Get(s, "/test", func(ctx *ContextNoBody) (string, error) {
		return "test", nil
	}).
		WithQueryParam("my-param", "my description")

	require.Equal(t, "my description", route.operation.Parameters.GetByInAndName("query", "my-param").Description)
}

func BenchmarkRequest(b *testing.B) {
	type Resp struct {
		Name string `json:"name"`
	}

	b.Run("fuego server and fuego post", func(b *testing.B) {
		s := NewServer()
		Post(s, "/test", func(c Ctx[MyStruct]) (Resp, error) {
			body, err := c.Body()
			if err != nil {
				return Resp{}, err
			}

			return Resp{Name: body.B}, nil
		})

		for i := 0; i < b.N; i++ {
			r := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"b":"M. John","c":3}`))
			w := httptest.NewRecorder()

			s.Mux.ServeHTTP(w, r)

			if w.Code != http.StatusOK || w.Body.String() != crlf(`{"name":"M. John"}`) {
				b.Fail()
			}
		}
	})

	b.Run("fuego server and std post", func(b *testing.B) {
		s := NewServer()
		PostStd(s, "/test", func(w http.ResponseWriter, r *http.Request) {
			var body MyStruct
			err := json.NewDecoder(r.Body).Decode(&body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			resp := Resp{
				Name: body.B,
			}
			err = json.NewEncoder(w).Encode(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		})

		for i := 0; i < b.N; i++ {
			r := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"b":"M. John","c":3}`))
			w := httptest.NewRecorder()

			s.Mux.ServeHTTP(w, r)

			if w.Code != http.StatusOK || w.Body.String() != crlf(`{"name":"M. John"}`) {
				b.Fail()
			}
		}
	})

	b.Run("std server and std post", func(b *testing.B) {
		mux := http.NewServeMux()
		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			var body MyStruct
			err := json.NewDecoder(r.Body).Decode(&body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			resp := Resp{
				Name: body.B,
			}
			err = json.NewEncoder(w).Encode(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		})

		for i := 0; i < b.N; i++ {
			r := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"b":"M. John","c":3}`))
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, r)

			if w.Code != http.StatusOK || w.Body.String() != crlf(`{"name":"M. John"}`) {
				b.Fail()
			}
		}
	})
}

func TestPerRouteMiddleware(t *testing.T) {
	s := NewServer()

	Get(s, "/withMiddleware", func(ctx *ContextNoBody) (string, error) {
		return "withmiddleware", nil
	}, dummyMiddleware)

	Get(s, "/withoutMiddleware", func(ctx *ContextNoBody) (string, error) {
		return "withoutmiddleware", nil
	})

	t.Run("withMiddleware", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/withMiddleware", nil)

		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, w.Body.String(), "\"withmiddleware\"\n")
		require.Equal(t, "response", w.Header().Get("X-Test-Response"))
	})

	t.Run("withoutMiddleware", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/withoutMiddleware", nil)

		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, w.Body.String(), "\"withoutmiddleware\"\n")
		require.Equal(t, "", w.Header().Get("X-Test-Response"))
	})
}

func TestGroup(t *testing.T) {
	s := NewServer()

	group1 := Group(s, "/group")
	Get(group1, "/route1", func(ctx *ContextNoBody) (string, error) {
		return "route1", nil
	})

	group2 := Group(s, "/group2")
	Use(group2, dummyMiddleware) // middleware is scoped to the group
	Get(group2, "/route2", func(ctx *ContextNoBody) (string, error) {
		return "route2", nil
	})

	subGroup := Group(group1, "/sub")

	Get(subGroup, "/route3", func(ctx *ContextNoBody) (string, error) {
		return "route3", nil
	})

	t.Run("route1", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/group/route1", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, w.Body.String(), "\"route1\"\n")
		require.Equal(t, "", w.Header().Get("X-Test-Response"), "middleware is not set to this group")
	})

	t.Run("route2", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/group2/route2", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, w.Body.String(), "\"route2\"\n")
		require.Equal(t, "response", w.Header().Get("X-Test-Response"))
	})

	t.Run("route3", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/group/sub/route3", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, w.Body.String(), "\"route3\"\n")
		require.Equal(t, "", w.Header().Get("X-Test-Response"), "middleware is not inherited")
	})
}
