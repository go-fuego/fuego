package fuego

import (
	"errors"
	"html/template"
	"io"
	"log/slog"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
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
	s := NewServer(
		WithXML(),
	)
	Get(s, "/", controller)
	Get(s, "/error", controllerWithError)

	t.Run("response is XML", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		s.Mux.ServeHTTP(recorder, req)

		require.Equal(t, 200, recorder.Code)
		require.Equal(t, "application/xml", recorder.Header().Get("Content-Type"))
		require.Equal(t, "<TestStruct><Name>Ewen</Name><Age>23</Age></TestStruct>", recorder.Body.String())
	})

	t.Run("error response is XML", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/error", nil)

		s.Mux.ServeHTTP(recorder, req)

		require.Equal(t, 500, recorder.Code)
		require.Equal(t, "application/xml", recorder.Header().Get("Content-Type"))
		require.Equal(t, "<HTTPError><title>Internal Server Error</title><status>500</status></HTTPError>", recorder.Body.String())
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

func TestWithOpenAPIGenerator(t *testing.T) {
	type User struct {
		Id *string `json:"test_id,omitempty" description:"test_description" readOnly:"true"`
	}

	t.Run("with default generator", func(t *testing.T) {
		s := NewServer()

		Get(s, "/users/{id}", func(*ContextNoBody) (User, error) {
			return User{}, nil
		})

		document := s.OutputOpenAPISpec()
		require.NotNil(t, document.Components.Schemas["User"].Value.Properties["test_id"])
		require.False(t, document.Components.Schemas["User"].Value.Properties["test_id"].Value.ReadOnly)
		require.Equal(t, document.Components.Schemas["User"].Value.Properties["test_id"].Value.Description, "")
	})

	t.Run("with custom generator", func(t *testing.T) {
		s := NewServer(
			WithOpenAPIGenerator(
				openapi3gen.NewGenerator(
					openapi3gen.UseAllExportedFields(),
					openapi3gen.SchemaCustomizer(func(name string, t reflect.Type, tag reflect.StructTag, schema *openapi3.Schema) error {
						if v := tag.Get("readOnly"); v == "true" {
							schema.ReadOnly = true
						}

						if v := tag.Get("description"); v != "" {
							schema.Description = v
						}

						return nil
					}),
				),
			),
		)

		Get(s, "/users/{id}", func(*ContextNoBody) (User, error) {
			return User{}, nil
		})

		document := s.OutputOpenAPISpec()
		require.NotNil(t, document.Components.Schemas["User"].Value.Properties["test_id"])
		require.True(t, document.Components.Schemas["User"].Value.Properties["test_id"].Value.ReadOnly)
		require.Equal(t, document.Components.Schemas["User"].Value.Properties["test_id"].Value.Description, "test_description")
	})
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
