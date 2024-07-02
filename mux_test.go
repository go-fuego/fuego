package fuego

import (
	"encoding/json"
	"fmt"
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

// orderMiddleware sets the X-Test-Order Header on the request and
// X-Test-Response header on the response. It is
// used to test the order execution of our middleware
func orderMiddleware(s string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Add("X-Test-Order", s)
			w.Header().Set("X-Test-Response", "response")
			handler.ServeHTTP(w, r)
		})
	}
}

// TestUse is used to mainly test the ordering of middleware execution
func TestUse(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		s := NewServer()
		Use(s, orderMiddleware("First!"))
		Get(s, "/test", func(ctx *ContextNoBody) (string, error) {
			return "test", nil
		})

		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		r.Header.Set("X-Test-Order", "Start!")
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, []string{"Start!", "First!"}, r.Header["X-Test-Order"])
	})

	t.Run("multiple uses of Use", func(t *testing.T) {
		s := NewServer()
		Use(s, orderMiddleware("First!"))
		Use(s, orderMiddleware("Second!"))
		Get(s, "/test", func(ctx *ContextNoBody) (string, error) {
			return "test", nil
		})

		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		r.Header.Set("X-Test-Order", "Start!")
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, []string{"Start!", "First!", "Second!"}, r.Header["X-Test-Order"])
	})

	t.Run("variadic use of Use", func(t *testing.T) {
		s := NewServer()
		Use(s, orderMiddleware("First!"))
		Use(s, orderMiddleware("Second!"), orderMiddleware("Third!"))
		Get(s, "/test", func(ctx *ContextNoBody) (string, error) {
			return "test", nil
		})

		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		r.Header.Set("X-Test-Order", "Start!")
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, []string{"Start!", "First!", "Second!", "Third!"}, r.Header["X-Test-Order"])
	})

	t.Run("variadic use of Route Get", func(t *testing.T) {
		s := NewServer()
		Use(s, orderMiddleware("First!"))
		Use(s, orderMiddleware("Second!"), orderMiddleware("Third!"))
		Get(s, "/test", func(ctx *ContextNoBody) (string, error) {
			return "test", nil
		}, orderMiddleware("Fourth!"), orderMiddleware("Fifth!"))

		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		r.Header.Set("X-Test-Order", "Start!")
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, []string{"Start!", "First!", "Second!", "Third!", "Fourth!", "Fifth!"}, r.Header["X-Test-Order"])
	})

	t.Run("group middlewares", func(t *testing.T) {
		s := NewServer()
		Use(s, orderMiddleware("First!"))
		group := Group(s, "/group")
		Use(group, orderMiddleware("Second!"))
		Use(group, orderMiddleware("Third!"))
		Get(group, "/test", func(ctx *ContextNoBody) (string, error) {
			return "test", nil
		})

		r := httptest.NewRequest(http.MethodGet, "/group/test", nil)
		r.Header.Set("X-Test-Order", "Start!")
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, []string{"Start!", "First!", "Second!", "Third!"}, r.Header["X-Test-Order"])
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
		require.Equal(t, "test", w.Body.String())
	})

	t.Run("post", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/test", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, w.Code, http.StatusOK)
		require.Equal(t, "test", w.Body.String())
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
	require.Equal(t, "test", w.Body.String())
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
	require.Equal(t, "test", w.Body.String())
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
	require.Equal(t, "test", w.Body.String())
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
	require.Equal(t, "test", w.Body.String())
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
	require.Equal(t, "test", "test", w.Body.String())
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

func TestAllStd(t *testing.T) {
	s := NewServer()
	AllStd(s, "/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test successful"))
	})

	t.Run("get", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPatch, "/test", nil)

		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, w.Code, http.StatusOK)
		require.Equal(t, w.Body.String(), "test successful")
	})

	t.Run("post", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/test", nil)

		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, w.Code, http.StatusOK)
		require.Equal(t, w.Body.String(), "test successful")
	})
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

func TestGroupTagsOnRoute(t *testing.T) {
	t.Run("route tag inheritance", func(t *testing.T) {
		s := NewServer().
			Tags("my-server-tag")
		route := Get(s, "/path", func(ctx *ContextNoBody) (string, error) {
			return "test", nil
		})
		require.Equal(t, []string{"my-server-tag"}, route.Operation.Tags)
	})

	t.Run("route tag override", func(t *testing.T) {
		s := NewServer().
			Tags("my-server-tag")

		route := Get(s, "/path", func(ctx *ContextNoBody) (string, error) {
			return "test", nil
		}).Tags("my-route-tag")

		require.Equal(t, []string{"my-route-tag"}, route.Operation.Tags)
	})

	t.Run("route tag add", func(t *testing.T) {
		s := NewServer().
			Tags("my-server-tag")

		route := Get(s, "/path", func(ctx *ContextNoBody) (string, error) {
			return "test", nil
		}).AddTags("my-route-tag")

		require.Equal(t, []string{"my-server-tag", "my-route-tag"}, route.Operation.Tags)
	})

	t.Run("route tag removal", func(t *testing.T) {
		s := NewServer().
			Tags("my-server-tag")

		route := Get(s, "/path", func(ctx *ContextNoBody) (string, error) {
			return "test", nil
		}).AddTags("my-route-tag").RemoveTags("my-server-tag")

		require.Equal(t, []string{"my-route-tag"}, route.Operation.Tags)
	})
}

func TestHideOpenapiRoutes(t *testing.T) {
	t.Run("hide main server", func(t *testing.T) {
		s := NewServer()
		Get(s, "/not-hidden", func(ctx *ContextNoBody) (string, error) { return "", nil })
		s.Hide()
		Get(s, "/test", func(ctx *ContextNoBody) (string, error) { return "", nil })

		require.Equal(t, s.DisableOpenapi, true)
		require.True(t, s.OpenApiSpec.Paths.Find("/not-hidden") != nil)
		require.True(t, s.OpenApiSpec.Paths.Find("/test") == nil)
	})

	t.Run("hide group", func(t *testing.T) {
		s := NewServer()
		Get(s, "/not-hidden", func(ctx *ContextNoBody) (string, error) { return "", nil })

		g := Group(s, "/group").Hide()
		Get(g, "/test", func(ctx *ContextNoBody) (string, error) { return "", nil })

		require.Equal(t, g.DisableOpenapi, true)
		require.True(t, s.OpenApiSpec.Paths.Find("/not-hidden") != nil)
		require.True(t, s.OpenApiSpec.Paths.Find("/group/test") == nil)
	})

	t.Run("hide group but not other group", func(t *testing.T) {
		s := NewServer()
		g := Group(s, "/group").Hide()
		Get(g, "/test", func(ctx *ContextNoBody) (string, error) { return "test", nil })

		g2 := Group(s, "/group2")
		Get(g2, "/test", func(ctx *ContextNoBody) (string, error) { return "test", nil })

		require.Equal(t, true, g.DisableOpenapi)
		require.Equal(t, false, g2.DisableOpenapi)
		require.True(t, s.OpenApiSpec.Paths.Find("/group/test") == nil)
		require.True(t, s.OpenApiSpec.Paths.Find("/group2/test") != nil)
	})

	t.Run("hide group but show sub group", func(t *testing.T) {
		s := NewServer()
		g := Group(s, "/group").Hide()
		Get(g, "/test", func(ctx *ContextNoBody) (string, error) { return "test", nil })

		g2 := Group(g, "/sub").Show()
		Get(g2, "/test", func(ctx *ContextNoBody) (string, error) { return "test", nil })

		require.Equal(t, true, g.DisableOpenapi)
		require.True(t, s.OpenApiSpec.Paths.Find("/group/test") == nil)
		require.True(t, s.OpenApiSpec.Paths.Find("/group/sub/test") != nil)
	})
}

func BenchmarkRequest(b *testing.B) {
	type Resp struct {
		Name string `json:"name"`
	}

	b.Run("fuego server and fuego post", func(b *testing.B) {
		s := NewServer()
		Post(s, "/test", func(c *ContextWithBody[MyStruct]) (Resp, error) {
			body, err := c.Body()
			if err != nil {
				return Resp{}, err
			}

			return Resp{Name: body.B}, nil
		})

		for range b.N {
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

		for range b.N {
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

		for range b.N {
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

		require.Equal(t, "withmiddleware", w.Body.String())
		require.Equal(t, "response", w.Header().Get("X-Test-Response"))
	})

	t.Run("withoutMiddleware", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/withoutMiddleware", nil)

		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, "withoutmiddleware", w.Body.String())
		require.Equal(t, "", w.Header().Get("X-Test-Response"))
	})
}

func TestGroup(t *testing.T) {
	s := NewServer()

	main := Group(s, "/")
	Use(main, dummyMiddleware) // middleware is scoped to the group
	Get(main, "/main", func(ctx *ContextNoBody) (string, error) {
		return "main", nil
	})

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

		require.Equal(t, "route1", w.Body.String())
		require.Equal(t, "", w.Header().Get("X-Test-Response"), "middleware is not set to this group")
	})

	t.Run("route2", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/group2/route2", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, "route2", w.Body.String())
		require.Equal(t, "response", w.Header().Get("X-Test-Response"))
	})

	t.Run("route3", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/group/sub/route3", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, "route3", w.Body.String())
		require.Equal(t, "", w.Header().Get("X-Test-Response"), "middleware is not inherited")
	})

	t.Run("main group", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/main", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, "main", w.Body.String())
		require.Equal(t, "response", w.Header().Get("X-Test-Response"), "middleware is not set to this group")
	})

	t.Run("group path can end with a slash (but with a warning)", func(t *testing.T) {
		s := NewServer()
		g := Group(s, "/slash/")
		require.Equal(t, "/slash/", g.basePath)
	})
}

func TestGroupTags(t *testing.T) {
	t.Run("inherit tags", func(t *testing.T) {
		s := NewServer().
			Tags("my-server-tag")
		group := Group(s, "/slash")

		require.Equal(t, []string{"my-server-tag"}, group.tags)
	})
	t.Run("override parent tags", func(t *testing.T) {
		s := NewServer().
			Tags("my-server-tag")
		group := Group(s, "/slash").
			Tags("my-group-tag")

		require.Equal(t, []string{"my-group-tag"}, group.tags)
	})
	t.Run("add child group tag", func(t *testing.T) {
		s := NewServer().
			Tags("my-server-tag")
		group := Group(s, "/slash").
			AddTags("my-group-tag")

		require.Equal(t, []string{"my-server-tag", "my-group-tag"}, group.tags)
	})
	t.Run("remove server tag", func(t *testing.T) {
		s := NewServer().
			Tags("my-server-tag", "my-other-server-tag")
		group := Group(s, "/slash").
			RemoveTags("my-server-tag")

		require.Equal(t, []string{"my-other-server-tag"}, group.tags)
	})
	t.Run("multiple groups inheritance", func(t *testing.T) {
		s := NewServer().
			Tags("my-server-tag")
		group := Group(s, "/slash").
			AddTags("my-group-tag")
		childGroup := Group(group, "/slash").
			AddTags("my-childGroup-tag")

		require.Equal(t, []string{"my-server-tag", "my-group-tag", "my-childGroup-tag"}, childGroup.tags)
	})
}

func ExampleContextNoBody_SetCookie() {
	s := NewServer()
	Get(s, "/test", func(c *ContextNoBody) (string, error) {
		c.SetCookie(http.Cookie{
			Name:  "name",
			Value: "value",
		})
		return "test", nil
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	s.Mux.ServeHTTP(w, r)

	fmt.Println(w.Result().Cookies()[0].Name)
	fmt.Println(w.Result().Cookies()[0].Value)

	// Output:
	// name
	// value
}

func ExampleContextNoBody_SetHeader() {
	s := NewServer()
	Get(s, "/test", func(c *ContextNoBody) (string, error) {
		c.SetHeader("X-Test", "test")
		return "test", nil
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	s.Mux.ServeHTTP(w, r)

	fmt.Println(w.Header().Get("X-Test"))

	// Output:
	// test
}

func wrappedFunc(custom string) func(string) string {
	return func(s string) string {
		return s + custom
	}
}

func TestNameFromNamespace(t *testing.T) {
	testCases := []struct {
		name string

		opts           []func(string) string
		route          Route[any, any]
		expectedOutput string
	}{
		{
			name: "base",

			route: Route[any, any]{
				FullName: "pkg.test.MyFunc1",
			},
			expectedOutput: "MyFunc1",
		},
		{
			name: "with camelToHuman",

			route: Route[any, any]{
				FullName: "pkg.test.MyFunc1",
			},
			opts: []func(string) string{
				camelToHuman,
			},
			expectedOutput: "my func1",
		},
		{
			name: "with inline opt",

			route: Route[any, any]{
				FullName: "pkg.test.MyFunc1",
			},
			opts: []func(string) string{
				camelToHuman,
				func(s string) string {
					return s + " foo"
				},
			},
			expectedOutput: "my func1 foo",
		},
		{
			name: "with wrapped func",

			route: Route[any, any]{
				FullName: "pkg.test.MyFunc1",
			},
			opts: []func(string) string{
				wrappedFunc("Foo"),
				camelToHuman,
			},
			expectedOutput: "my func1 foo",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.route.NameFromNamespace(tc.opts...)
			require.Equal(t, tc.expectedOutput, actual)
		})
	}
}

func BenchmarkCamelToHuman(b *testing.B) {
	b.Run("camelToHuman", func(b *testing.B) {
		for range b.N {
			camelToHuman("listAllRecipes")
		}
	})
}

func TestCamelToHuman(t *testing.T) {
	testCases := []struct {
		input  string
		output string
	}{
		{"listAllRecipes", "list all recipes"},
		{"get5Recipes", "get5 recipes"},
		{"getHTTP", "get h t t p"},
		{"getHTTP2", "get h t t p2"},
		{"getHTTP2Server", "get h t t p2 server"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			require.Equal(t, tc.output, camelToHuman(tc.input))
		})
	}
}
