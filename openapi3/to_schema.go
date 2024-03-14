package openapi3

import (
	"reflect"
	"strings"
	"time"
)

// ToSchema converts any Go type to an OpenAPI Schema
func ToSchema(v any) *Schema {
	if v == nil {
		return nil
	}

	s := Schema{
		Type:       Object,
		Properties: make(map[string]Schema),
	}

	value := reflect.ValueOf(v)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if _, isTime := v.(time.Time); isTime {
		s.Type = "string"
		s.Format = "date-time"
		s.Example = time.RFC3339
		return &Schema{
			Type:    "string",
			Format:  "date-time",
			Example: time.RFC3339,
		}
	}

	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		s.Type = Array
		itemType := value.Type().Elem()
		if itemType.Kind() == reflect.Ptr {
			itemType = itemType.Elem()
		}
		one := reflect.New(itemType)
		s.Items = ToSchema(one.Interface())
		s.Example = "[]"
	case reflect.Struct:
		for i := range value.NumField() {
			structField := value.Type().Field(i)

			fieldName := fieldName(structField)

			fieldValue := value.Field(i)

			fieldSchema := *ToSchema(fieldValue.Interface())

			// Parse struct tags
			if strings.Contains(structField.Tag.Get("validate"), "required") {
				s.Required = append(s.Required, fieldName)
			}
			parseValidate(&fieldSchema, structField.Tag.Get("validate"))
			fieldSchema.Example = structField.Tag.Get("example")
			fieldSchema.Format = structField.Tag.Get("format")
			s.Properties[fieldName] = fieldSchema
		}
	default:
		s.Type = kindToType(value.Kind())
	}

	return &s
}

func fieldName(s reflect.StructField) string {
	jsonTags := strings.Split(s.Tag.Get("json"), ",")
	if len(jsonTags) > 0 && jsonTags[0] != "" {
		return jsonTags[0]
	}
	return s.Name
}

type OpenAPIType string

const (
	Invalid OpenAPIType = ""
	String  OpenAPIType = "string"
	Integer OpenAPIType = "integer"
	Number  OpenAPIType = "number"
	Boolean OpenAPIType = "boolean"
	Array   OpenAPIType = "array"
	Object  OpenAPIType = "object"
)

func kindToType(kind reflect.Kind) OpenAPIType {
	switch kind {
	case reflect.String:
		return String
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Integer
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return Integer
	case reflect.Float32, reflect.Float64:
		return Number
	case reflect.Bool:
		return Boolean
	case reflect.Slice, reflect.Array:
		return Array
	case reflect.Struct:
		return Object
	default:
		return Invalid
	}
}
