package openapi3

import (
	"reflect"
	"strings"
	"time"
)

func NewDocument() Document {
	return Document{
		OpenAPI: "3.1.0",
		Info: Info{
			Title:       "OpenAPI",
			Description: "OpenAPI",
		},
		Paths:      make(Paths),
		Components: NewComponents(),
	}
}

type Document struct {
	OpenAPI string `json:"openapi" yaml:"openapi"`
	Info    Info   `json:"info" yaml:"info"`

	Paths      Paths      `json:"paths" yaml:"paths"`
	Components Components `json:"components" yaml:"components"`
}

type Paths map[string]map[string]*Operation

func (p Paths) AddPath(path string, method string, pathItem *Operation) {
	if p[path] == nil {
		p[path] = make(map[string]*Operation)
	}
	p[path][method] = pathItem
}

type Info struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Version     string `json:"version" yaml:"version"`
}

type Schema struct {
	Type       string            `json:"type" yaml:"type"`
	Format     string            `json:"format,omitempty" yaml:"format,omitempty"`
	Required   []string          `json:"required,omitempty" yaml:"required,omitempty"`
	Example    string            `json:"example,omitempty" yaml:"example,omitempty"`
	Properties map[string]Schema `json:"properties,omitempty" yaml:"properties,omitempty"`
	Items      *Schema           `json:"items,omitempty" yaml:"items,omitempty"`
}

// ToSchema converts any Go type to an OpenAPI Schema
func ToSchema(v any) *Schema {
	if v == nil {
		return nil
	}

	s := Schema{
		Type:       "object",
		Properties: make(map[string]Schema),
	}

	value := reflect.ValueOf(v)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() == reflect.Slice {
		s.Type = "array"
		one := reflect.New(value.Type().Elem())
		s.Items = ToSchema(one.Interface())
	}

	if _, isTime := value.Interface().(time.Time); isTime {
		s.Type = "string"
		s.Format = "date-time"
		s.Example = value.Interface().(time.Time).Format(time.RFC3339)
		return &s
	}

	if value.Kind() == reflect.Struct {
		// Iterate on fields with reflect
		for i := range value.NumField() {
			field := value.Field(i)
			fieldType := value.Type().Field(i)

			// If the field is a struct, we need to dive into it
			if field.Kind() == reflect.Struct {
				fieldName := fieldType.Tag.Get("json")
				if fieldName == "" {
					fieldName = fieldType.Name
				}
				s.Properties[fieldName] = *ToSchema(field.Interface())
			} else {
				// If the field is a basic type, we can just add it to the properties
				fieldTypeType := fieldType.Type.Name()
				format := fieldType.Tag.Get("format")
				if strings.Contains(fieldTypeType, "int") {
					fieldTypeType = "integer"
					if format != "" {
						format = fieldType.Type.Name()
					}
				} else if fieldTypeType == "bool" {
					fieldTypeType = "boolean"
				}
				fieldName := fieldType.Tag.Get("json")
				if fieldName == "" {
					fieldName = fieldType.Name
				}
				if strings.Contains(fieldType.Tag.Get("validate"), "required") {
					s.Required = append(s.Required, fieldName)
				}
				s.Properties[fieldName] = Schema{
					Type:    fieldTypeType,
					Example: fieldType.Tag.Get("example"),
					Format:  format,
				}
			}
		}
	}

	if !(value.Kind() == reflect.Struct || value.Kind() == reflect.Slice) {
		s.Type = value.Kind().String()
		if strings.Contains(s.Type, "int") {
			s.Type = "integer"
		} else if s.Type == "bool" {
			s.Type = "boolean"
		}
	}

	return &s
}

type Parameter struct {
	Name        string `json:"name" yaml:"name"`
	In          string `json:"in" yaml:"in"`
	Description string `json:"description" yaml:"description"`
	Required    bool   `json:"required" yaml:"required"`
	Schema      Schema `json:"schema" yaml:"schema"`
	Example     string `json:"example,omitempty" yaml:"example"`
}

type MimeType string

type Response struct {
	Description string                    `json:"description" yaml:"description"`
	Content     map[MimeType]SchemaObject `json:"content" yaml:"content"`
}

type SchemaObject struct {
	Schema *Schema `json:"schema" yaml:"schema"`
}

type Operation struct {
	OperationID string               `json:"operationId" yaml:"operationId"`
	Summary     string               `json:"summary" yaml:"summary"`
	Description string               `json:"description" yaml:"description"`
	Deprecated  bool                 `json:"deprecated,omitempty" yaml:"deprecated"`
	RequestBody *RequestBody         `json:"requestBody,omitempty" yaml:"requestBody"`
	Parameters  []*Parameter         `json:"parameters,omitempty" yaml:"parameters"`
	Tags        []string             `json:"tags,omitempty" yaml:"tags"`
	Responses   map[string]*Response `json:"responses,omitempty" yaml:"responses"`
}

type RequestBody struct {
	Required bool                      `json:"required" yaml:"required"`
	Content  map[MimeType]SchemaObject `json:"content" yaml:"content"`
}

type Parameters struct {
	Name string `json:"name" yaml:"name"`
	In   string `json:"in" yaml:"in"`
}

func NewComponents() Components {
	return Components{
		Schemas:       make(map[string]*Schema),
		RequestBodies: make(map[string]*RequestBody),
	}
}

type Components struct {
	Schemas       map[string]*Schema      `json:"schemas" yaml:"schemas"`
	RequestBodies map[string]*RequestBody `json:"requestBodies" yaml:"requestBodies"`
}
