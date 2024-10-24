package fuego

import (
	"errors"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func controller(c *ContextNoBody) (testStruct, error) {
	return testStruct{Name: "Ewen", Age: 23}, nil
}

func controllerWithError(c *ContextNoBody) (testStruct, error) {
	return testStruct{}, errors.New("error")
}

func TestNewServer(t *testing.T) {
	s := NewServer()

	t.Run("can register controller", func(t *testing.T) {
		Get(s, "/", controller)

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		s.Mux.ServeHTTP(recorder, req)

		require.Equal(t, 200, recorder.Code)
	})
}

func TestWithXML(t *testing.T) {
	s := NewServer()
	Get(s, "/", controller)
	Get(s, "/error", controllerWithError)

	t.Run("response is XML", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Accept", "application/xml")

		s.Mux.ServeHTTP(recorder, req)

		require.Equal(t, 200, recorder.Code)
		require.Equal(t, "application/xml", recorder.Header().Get("Content-Type"))
		require.Equal(t, "<TestStruct><Name>Ewen</Name><Age>23</Age></TestStruct>", recorder.Body.String())
	})

	t.Run("error response is XML", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/error", nil)
		req.Header.Set("Accept", "application/xml")

		s.Mux.ServeHTTP(recorder, req)

		require.Equal(t, 500, recorder.Code)
		require.Equal(t, "<HTTPError><title>Internal Server Error</title><status>500</status></HTTPError>", recorder.Body.String())
		require.Equal(t, "application/xml", recorder.Header().Get("Content-Type"))
	})
}

func TestWithOpenAPIConfig(t *testing.T) {
	t.Run("with default values", func(t *testing.T) {
		s := NewServer(
			WithOpenAPIConfig(OpenAPIConfig{}),
		)

		require.Equal(t, "/swagger", s.OpenAPIConfig.SwaggerUrl)
		require.Equal(t, "/swagger/openapi.json", s.OpenAPIConfig.JsonUrl)
		require.Equal(t, "doc/openapi.json", s.OpenAPIConfig.JsonFilePath)
		require.False(t, s.OpenAPIConfig.PrettyFormatJson)
	})

	t.Run("with custom values", func(t *testing.T) {
		s := NewServer(
			WithOpenAPIConfig(OpenAPIConfig{
				SwaggerUrl:       "/api",
				JsonUrl:          "/api/openapi.json",
				JsonFilePath:     "openapi.json",
				DisableSwagger:   true,
				DisableLocalSave: true,
				PrettyFormatJson: true,
			}),
		)

		require.Equal(t, "/api", s.OpenAPIConfig.SwaggerUrl)
		require.Equal(t, "/api/openapi.json", s.OpenAPIConfig.JsonUrl)
		require.Equal(t, "openapi.json", s.OpenAPIConfig.JsonFilePath)
		require.True(t, s.OpenAPIConfig.DisableSwagger)
		require.True(t, s.OpenAPIConfig.DisableLocalSave)
		require.True(t, s.OpenAPIConfig.PrettyFormatJson)
	})

	t.Run("with invalid local path values", func(t *testing.T) {
		t.Run("with invalid path", func(t *testing.T) {
			NewServer(
				WithOpenAPIConfig(OpenAPIConfig{
					JsonFilePath: "path/to/jsonSpec",
					SwaggerUrl:   "p   i",
					JsonUrl:      "pi/op  enapi.json",
				}),
			)
		})
		t.Run("with invalid url", func(t *testing.T) {
			NewServer(
				WithOpenAPIConfig(OpenAPIConfig{
					JsonFilePath: "path/to/jsonSpec.json",
					JsonUrl:      "pi/op  enapi.json",
					SwaggerUrl:   "p   i",
				}),
			)
		})

		t.Run("with invalid url", func(t *testing.T) {
			NewServer(
				WithOpenAPIConfig(OpenAPIConfig{
					JsonFilePath: "path/to/jsonSpec.json",
					JsonUrl:      "/api/openapi.json",
					SwaggerUrl:   "invalid path",
				}),
			)
		})
	})
}

func TestWithBasePath(t *testing.T) {
	s := NewServer(
		WithBasePath("/api"),
	)

	require.Equal(t, "/api", s.basePath)
}

func TestWithMaxBodySize(t *testing.T) {
	s := NewServer(
		WithMaxBodySize(1024),
	)

	require.Equal(t, int64(1024), s.maxBodySize)
}

func TestWithAutoAuth(t *testing.T) {
	s := NewServer(
		WithAutoAuth(nil),
	)

	require.NotNil(t, s.autoAuth)
	require.True(t, s.autoAuth.Enabled)
	// The authoauth is tested in security_test.go,
	// this is just an option to enable it.
}

func TestWithTemplates(t *testing.T) {
	t.Run("with template FS", func(t *testing.T) {
		template := template.New("test")
		s := NewServer(
			WithTemplateFS(testdata),
			WithTemplates(template),
		)

		require.NotNil(t, s.template)
	})

	t.Run("without template FS", func(t *testing.T) {
		template := template.New("test")
		s := NewServer(
			WithTemplates(template),
		)

		require.NotNil(t, s.template)
	})
}

func TestWithLogHandler(t *testing.T) {
	handler := slog.NewTextHandler(io.Discard, nil)
	NewServer(
		WithLogHandler(handler),
	)
}

func TestWithValidator(t *testing.T) {
	type args struct {
		newValidator *validator.Validate
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name: "with custom validator",
			args: args{
				newValidator: validator.New(),
			},
		},
		{
			name: "no validator provided",
			args: args{
				newValidator: nil,
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if tt.wantPanic {
					assert.Panics(
						t, func() { WithValidator(tt.args.newValidator) },
					)
				} else {
					NewServer(
						WithValidator(tt.args.newValidator),
					)
					assert.Equal(t, tt.args.newValidator, v)
				}
			},
		)
	}
}

func TestWithAddr(t *testing.T) {
	tests := []struct {
		name         string
		addr         string
		expectedAddr string
	}{
		{
			name:         "with custom addr, that addr is used",
			addr:         "localhost:8888",
			expectedAddr: "localhost:8888",
		},
		{
			name:         "no addr provided, default is used (9999)",
			addr:         "",
			expectedAddr: "localhost:9999",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				var opts []func(*Server)

				if tt.addr != "" {
					opts = append(opts, WithAddr(tt.addr))
				}

				s := NewServer(
					opts...,
				)
				require.Equal(t, tt.expectedAddr, s.Server.Addr)
			},
		)
	}
}

func TestWithPort(t *testing.T) {
	t.Run("with custom port, that port is used", func(t *testing.T) {
		s := NewServer(
			WithPort(8488),
		)
		require.Equal(t, "localhost:8488", s.Server.Addr)
	})

	t.Run("no port provided, default is used (9999)", func(t *testing.T) {
		s := NewServer()
		require.Equal(t, "localhost:9999", s.Server.Addr)
	})
}

func TestWithoutStartupMessages(t *testing.T) {
	s := NewServer(
		WithoutStartupMessages(),
	)

	require.True(t, s.disableStartupMessages)
}

func TestWithoutAutoGroupTags(t *testing.T) {
	s := NewServer(
		WithoutAutoGroupTags(),
	)

	require.True(t, s.disableAutoGroupTags)

	group := Group(s, "/api")
	Get(group, "/test", controller)

	document := s.OutputOpenAPISpec()
	require.NotNil(t, document)
	require.Nil(t, document.Paths.Find("/api/test").Get.Tags)
}

func TestServerTags(t *testing.T) {
	t.Run("set tags", func(t *testing.T) {
		s := NewServer().
			Tags("my-server-tag")

		require.Equal(t, s.tags, []string{"my-server-tag"})
	})

	t.Run("add tags", func(t *testing.T) {
		s := NewServer().
			AddTags("my-server-tag").
			AddTags("my-other-server-tag")

		require.Equal(t, s.tags, []string{"my-server-tag", "my-other-server-tag"})
	})

	t.Run("remove tags", func(t *testing.T) {
		s := NewServer().
			Tags("my-server-tag").
			AddTags("my-other-server-tag").
			RemoveTags("my-other-server-tag")

		require.Equal(t, s.tags, []string{"my-server-tag"})
	})

	t.Run("inherit tags from group, replace", func(t *testing.T) {
		s := NewServer().
			Tags("my-server-tag")

		group := Group(s, "/api").
			Tags("my-group-tag")

		require.Equal(t, group.tags, []string{"my-group-tag"})

		subGroup := Group(group, "/users").
			Tags("my-sub-group-tag")

		require.Equal(t, subGroup.tags, []string{"my-sub-group-tag"})
	})

	t.Run("inherit tags from group, add", func(t *testing.T) {
		s := NewServer().
			Tags("my-server-tag")

		group := Group(s, "/api").
			AddTags("my-group-tag")

		require.Equal(t, group.tags, []string{"my-server-tag", "my-group-tag"})

		subGroup := Group(group, "/users").
			AddTags("my-sub-group-tag")

		require.Equal(t, subGroup.tags, []string{"my-server-tag", "my-group-tag", "my-sub-group-tag"})
	})

	t.Run("inherit tags from group, remove", func(t *testing.T) {
		s := NewServer().
			Tags("my-server-tag")

		group := Group(s, "/api").
			AddTags("my-group-tag")

		require.Equal(t, group.tags, []string{"my-server-tag", "my-group-tag"})

		siblingGroup := Group(s, "/api2").
			AddTags("my-sibling-group-tag")

		require.Equal(t, siblingGroup.tags, []string{"my-server-tag", "my-sibling-group-tag"})

		subGroup := Group(group, "/users").
			RemoveTags("my-group-tag")

		require.Equal(t, subGroup.tags, []string{"my-server-tag"})
	})
}

type ReqBody struct {
	A string
	B int
}

type Resp struct {
	Message string `json:"message"`
}

func dummyController(_ *ContextWithBody[ReqBody]) (Resp, error) {
	return Resp{Message: "hello world"}, nil
}

func TestWithRequestContentType(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		s := NewServer()
		require.Nil(t, s.acceptedContentTypes)
	})

	t.Run("input", func(t *testing.T) {
		arr := []string{"application/json", "application/xml"}
		s := NewServer(WithRequestContentType("application/json", "application/xml"))
		require.ElementsMatch(t, arr, s.acceptedContentTypes)
	})

	t.Run("ensure applied to route", func(t *testing.T) {
		s := NewServer(WithRequestContentType("application/json", "application/xml"))
		route := Post(s, "/test", dummyController)

		content := route.Operation.RequestBody.Value.Content
		require.NotNil(t, content.Get("application/json"))
		require.NotNil(t, content.Get("application/xml"))
		require.Equal(t, "#/components/schemas/ReqBody", content.Get("application/json").Schema.Ref)
		require.Equal(t, "#/components/schemas/ReqBody", content.Get("application/xml").Schema.Ref)
		_, ok := s.OpenApiSpec.Components.RequestBodies["ReqBody"]
		require.False(t, ok)
	})
}

func TestCustomSerialization(t *testing.T) {
	s := NewServer(
		WithSerializer(func(w http.ResponseWriter, r *http.Request, a any) error {
			w.WriteHeader(202)
			_, err := w.Write([]byte("custom serialization"))
			return err
		}),
	)

	Get(s, "/", func(c *ContextNoBody) (ans, error) {
		return ans{Ans: "Hello World"}, nil
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	require.Equal(t, 202, w.Code)
	require.Equal(t, "custom serialization", w.Body.String())
}

func TestGroupParams(t *testing.T) {
	s := NewServer()
	group := Group(s, "/api").Param("X-Test-Header", "test-value", OpenAPIParamOption{Example: "example", Required: true})

	Get(s, "/", controller)
	Get(group, "/test", controller)
	Get(group, "/test2", controller)

	require.Equal(t, group.params, []OpenAPIParam{{Name: "X-Test-Header", Description: "test-value", OpenAPIParamOption: OpenAPIParamOption{Required: true, Example: "example", Type: ""}}})

	document := s.OutputOpenAPISpec()
	require.Len(t, document.Paths.Find("/").Get.Parameters, 1)
	require.Equal(t, document.Paths.Find("/").Get.Parameters[0].Value.Name, "Accept")
	require.Equal(t, document.Paths.Find("/api/test").Get.Parameters[0].Value.Name, "Accept")
	require.Equal(t, document.Paths.Find("/api/test").Get.Parameters[1].Value.Name, "X-Test-Header")
	require.Equal(t, document.Paths.Find("/api/test2").Get.Parameters[0].Value.Name, "Accept")
	require.Equal(t, document.Paths.Find("/api/test2").Get.Parameters[1].Value.Name, "X-Test-Header")
}

func TestGroupHeaderParams(t *testing.T) {
	s := NewServer()
	group := Group(s, "/api").Header("X-Test-Header", "test-value")

	Get(group, "/test", controller)

	require.Equal(t, group.params, []OpenAPIParam{{Name: "X-Test-Header", Description: "test-value", OpenAPIParamOption: OpenAPIParamOption{Required: false, Example: "", Type: HeaderParamType}}})

	document := s.OutputOpenAPISpec()
	require.Equal(t, document.Paths.Find("/api/test").Get.Parameters[0].Value.Name, "Accept")
	require.Equal(t, document.Paths.Find("/api/test").Get.Parameters[1].Value.Name, "X-Test-Header")
}

func TestGroupCookieParams(t *testing.T) {
	s := NewServer()
	group := Group(s, "/api").Cookie("X-Test-Cookie", "test-value")

	Get(group, "/test", controller)

	require.Equal(t, group.params, []OpenAPIParam{{Name: "X-Test-Cookie", Description: "test-value", OpenAPIParamOption: OpenAPIParamOption{Required: false, Example: "", Type: CookieParamType}}})

	document := s.OutputOpenAPISpec()
	require.Equal(t, document.Paths.Find("/api/test").Get.Parameters[0].Value.Name, "Accept")
	require.Equal(t, document.Paths.Find("/api/test").Get.Parameters[1].Value.Name, "X-Test-Cookie")
}

func TestGroupQueryParam(t *testing.T) {
	s := NewServer()
	group := Group(s, "/api").Query("X-Test-Query", "test-value")

	Get(group, "/test", controller)

	require.Equal(t, group.params, []OpenAPIParam{{Name: "X-Test-Query", Description: "test-value", OpenAPIParamOption: OpenAPIParamOption{Required: false, Example: "", Type: QueryParamType}}})

	document := s.OutputOpenAPISpec()
	require.Equal(t, document.Paths.Find("/api/test").Get.Parameters[0].Value.Name, "Accept")
	require.Equal(t, document.Paths.Find("/api/test").Get.Parameters[1].Value.Name, "X-Test-Query")
}

func TestGroupParamsInChildGroup(t *testing.T) {
	s := NewServer()
	group := Group(s, "/api").Param("X-Test-Header", "test-value")

	subGroup := Group(group, "/users")

	expectedParams := []OpenAPIParam{{Name: "X-Test-Header", Description: "test-value", OpenAPIParamOption: OpenAPIParamOption{Required: false, Example: "", Type: ""}}}
	require.Equal(t, expectedParams, group.params)
	require.Equal(t, expectedParams, subGroup.params)
}

func TestGroupParamsNotInParentGroup(t *testing.T) {
	s := NewServer()
	parentGroup := Group(s, "/api")
	group := Group(parentGroup, "/users").Param("X-Test-Header", "test-value")

	expectedParams := []OpenAPIParam{{Name: "X-Test-Header", Description: "test-value", OpenAPIParamOption: OpenAPIParamOption{Required: false, Example: "", Type: ""}}}
	require.Nil(t, parentGroup.params)
	require.Equal(t, expectedParams, group.params)
}

func TestGroupParamsNotInSiblingGroup(t *testing.T) {
	s := NewServer()
	group := Group(s, "/api").Param("X-Test-Header", "test-value")
	siblingGroup := Group(s, "/api2")

	expectedParams := []OpenAPIParam{{Name: "X-Test-Header", Description: "test-value", OpenAPIParamOption: OpenAPIParamOption{Required: false, Example: "", Type: ""}}}
	require.Equal(t, expectedParams, group.params)
	require.Nil(t, siblingGroup.params)
}

func TestGroupParamsInMainServerInstance(t *testing.T) {
	s := NewServer().Param("X-Test-Header", "test-value")

	expectedParams := []OpenAPIParam{{Name: "X-Test-Header", Description: "test-value", OpenAPIParamOption: OpenAPIParamOption{Required: false, Example: "", Type: ""}}}
	require.Equal(t, expectedParams, s.params)
}

func TestHideGroupAfterGroupParam(t *testing.T) {
	s := NewServer()
	group := Group(s, "/api").Param("X-Test-Header", "test-value").Hide()

	Get(group, "/test", controller)

	document := s.OutputOpenAPISpec()
	require.Nil(t, document.Paths.Find("/api/test"))
}
