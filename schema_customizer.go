package fuego

import (
	"log/slog"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// parseValidate parses the values of the validate tag
// It adds the following struct tags (tag => OpenAPI schema field):
// - validate:
//   - min=1 => min=1 (for integers)
//   - min=1 => minLength=1 (for strings)
//   - max=100 => max=100 (for integers)
//   - max=100 => maxLength=100 (for strings)
//   - required => (*)
//
// (*) The value of required is set in "determineRequired" function.
// The "required" value is set using the "x-fuego-required-marker" temporary extension value.
func parseValidate(tag reflect.StructTag, schema *openapi3.Schema) {
	validateTag, ok := tag.Lookup("validate")
	if !ok {
		return
	}

	validateTags := strings.Split(validateTag, ",")
	for _, validateTag := range validateTags {
		// Mark what fields are needed for required. This will be removed later
		if validateTag == "required" {
			if schema.Extensions == nil {
				schema.Extensions = make(map[string]interface{})
			}
			schema.Extensions["x-fuego-required-marker"] = true
		}
		if strings.HasPrefix(validateTag, "min=") {
			minValue, err := strconv.Atoi(strings.Split(validateTag, "=")[1])
			if err != nil {
				slog.Warn("Min might be incorrect (should be integer)", "error", err)
			}

			if schema.Type.Is(openapi3.TypeInteger) {
				minPtr := float64(minValue)
				schema.Min = &minPtr
			} else if schema.Type.Is(openapi3.TypeString) {
				//nolint:gosec // disable G115
				schema.MinLength = uint64(minValue)
			}
		}
		if strings.HasPrefix(validateTag, "max=") {
			maxValue, err := strconv.Atoi(strings.Split(validateTag, "=")[1])
			if err != nil {
				slog.Warn("Max might be incorrect (should be integer)", "error", err)
			}
			if schema.Type.Is(openapi3.TypeInteger) {
				maxPtr := float64(maxValue)
				schema.Max = &maxPtr
			} else if schema.Type.Is(openapi3.TypeString) {
				//nolint:gosec // disable G115
				maxPtr := uint64(maxValue)
				schema.MaxLength = &maxPtr
			}
		}
	}
}

// determineRequired takes a struct type and determines which properties are
// required, as determined by checking the extension marker: "x-fuego-required-marker".
// This marker is set in "parseValidate" and is removed in this function.
func determineRequired(t reflect.Type, schema *openapi3.Schema) {
	if t.Kind() != reflect.Struct {
		return
	}

	var requiredFields []string

	for name, prop := range schema.Properties {
		// Check for our marker and remove it
		if val, ok := prop.Value.Extensions["x-fuego-required-marker"]; ok {
			if valBool, ok := val.(bool); ok && valBool {
				requiredFields = append(requiredFields, name)
			}
			delete(prop.Value.Extensions, "x-fuego-required-marker")
		}
	}

	sort.Strings(requiredFields)

	schema.Required = requiredFields
}

// parseOmitEmpty parses the "json" tag for the "omitempty" value.
// If it exists, the schema for the property sets "nullable" to true.
func parseOmitEmpty(tag reflect.StructTag, schema *openapi3.Schema) {
	jsonTag, ok := tag.Lookup("json")
	if ok && strings.Contains(jsonTag, ",omitempty") {
		schema.Nullable = true
	}
}

// parseExample parses the "example" tag and sets the schema example.
// If the example type does not match the type of the property, it will log a warning.
func parseExample(tag reflect.StructTag, schema *openapi3.Schema) {
	example, ok := tag.Lookup("example")
	if !ok {
		return
	}

	switch {
	case schema.Type.Is(openapi3.TypeInteger):
		exNum, err := strconv.Atoi(example)
		if err != nil {
			slog.Warn("Example might be incorrect (should be integer)", "error", err)
		}
		schema.Example = exNum
	case schema.Type.Is(openapi3.TypeNumber):
		exNum, err := strconv.ParseFloat(example, 64)
		if err != nil {
			slog.Warn("Example might be incorrect (should be floating point number)", "error", err)
		}
		schema.Example = exNum
	case schema.Type.Is(openapi3.TypeBoolean):
		exBool, err := strconv.ParseBool(example)
		if err != nil {
			slog.Warn("Example might be incorrect (should be boolean)", "error", err)
		}
		schema.Example = exBool
	default:
		schema.Example = example
	}
}

// parseDescriptions parses the "description" tag and adds it to the schema description.
func parseDescription(tag reflect.StructTag, schema *openapi3.Schema) {
	description, ok := tag.Lookup("description")
	if ok {
		schema.Description = description
	}
}

// SchemaCustomizer parses struct tags and modifies the schema using kin-openapi3gen's
// schema customization functionality.
// It adds the following struct tags (tag => OpenAPI schema field):
// - description => description
// - example => example
// - json => nullable (if contains omitempty)
// - validate:
//   - required => required
//   - min=1 => min=1 (for integers)
//   - min=1 => minLength=1 (for strings)
//   - max=100 => max=100 (for integers)
//   - max=100 => maxLength=100 (for strings)
func SchemaCustomizer(name string, t reflect.Type, tag reflect.StructTag, schema *openapi3.Schema) error {
	// Example
	parseExample(tag, schema)

	// Validation
	parseValidate(tag, schema)

	// Description
	parseDescription(tag, schema)

	// Omitempty
	parseOmitEmpty(tag, schema)

	// After we are done parsing tags, get the required tags
	determineRequired(t, schema)

	return nil
}
