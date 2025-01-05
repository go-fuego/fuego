package fuego

import (
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
)

var DefaultParsers = []MetadataParserEntry{
	{Name: "exampleParser", Parser: MetadataParserExample},
	{Name: "validationParser", Parser: MetadataParserValidation},
	{Name: "descriptionParser", Parser: MetadataParserDescription},
	{Name: "XMLParser", Parser: MetadataParserXML},
	{Name: "JSONParser", Parser: MetadataParserJSON},
}

type MetadataParserParams struct {
	Field      reflect.StructField    // Reflective information of a struct field
	FieldName  string                 // Name of the field
	Property   *openapi3.Schema       // OpenAPI schema of the field
	SchemaRef  *openapi3.SchemaRef    // Reference to the OpenAPI schema
	Additional map[string]interface{} // Additional metadata as key-value pairs
}

type MetadataParserFunction func(params MetadataParserParams)

type MetadataParserEntry struct {
	Name   string
	Parser MetadataParserFunction
}

type MetadataParsers struct {
	registeredParsers []MetadataParserEntry
	parserLock        sync.Mutex
	registeredNames   map[string]bool
}

// NewMetadataParsers returns a new MetadataParsers instance with empty registeredParsers
// slice and registeredNames map.
func NewMetadataParsers() *MetadataParsers {
	return &MetadataParsers{
		registeredParsers: []MetadataParserEntry{},
		registeredNames:   make(map[string]bool),
	}
}

// GetRegisteredParsers returns a slice of all currently registered metadata parsers.
// This method ensures thread-safe access to the registeredParsers slice by using a mutex lock.
func (mp *MetadataParsers) GetRegisteredParsers() []MetadataParserEntry {
	mp.parserLock.Lock()
	defer mp.parserLock.Unlock()

	return mp.registeredParsers
}

// Reset clears all registered metadata parsers from the slice and resets the
// registeredNames map as empty. This method ensures thread-safe access to the
// registeredParsers slice and registeredNames map by using a mutex lock.
func (mp *MetadataParsers) Reset() {
	mp.parserLock.Lock()
	defer mp.parserLock.Unlock()

	mp.registeredParsers = []MetadataParserEntry{}
	mp.registeredNames = make(map[string]bool)
}

// InitializeMetadataParsers initializes the metadata parsers with the given customParsers.
// It will append the customParsers to the existing DefaultParsers. If a parser with the
// same name is already registered, it will be skipped.
func (mp *MetadataParsers) InitializeMetadataParsers(customParsers []MetadataParserEntry) {
	mp.parserLock.Lock()
	defer mp.parserLock.Unlock()

	mp.registeredParsers = []MetadataParserEntry{}
	mp.registeredNames = make(map[string]bool)

	var parsersToRegister []MetadataParserEntry
	parsersToRegister = customParsers

	for _, entry := range parsersToRegister {
		if _, exists := mp.registeredNames[entry.Name]; exists {
			slog.Warn("Parser already registered", "name", entry.Name)
			continue
		}
		mp.registeredNames[entry.Name] = true
		mp.registeredParsers = append(mp.registeredParsers, entry)
	}
}

// RegisterMetadataParser registers a new metadata parser with the given name, parser
// function, and position string. The position string can be one of the following:
// "start", "end", "before", or "after". If the position is "before" or "after", you
// must also provide the name of the parser that you want to place the new parser
// before or after. If the relative parser is not found, an error is returned. If the
// parser is already registered, the function does nothing and returns nil. The
// function is thread-safe and ensures that the registeredParsers slice and
// registeredNames map are accessed in a thread-safe manner.
func (mp *MetadataParsers) RegisterMetadataParser(name string, parser MetadataParserFunction, position string, relativeTo string) error {
	validPositions := map[string]bool{"start": true, "end": true, "before": true, "after": true}
	if !validPositions[position] {
		return fmt.Errorf("Invalid position. Use 'start', 'end', 'before', or 'after'")
	}
	if (position == "before" || position == "after") && mp.findParserIndex(relativeTo) == -1 {
		return fmt.Errorf("Relative parser '%s' not found", relativeTo)
	}

	mp.parserLock.Lock()
	defer mp.parserLock.Unlock()

	for _, entry := range mp.registeredParsers {
		if entry.Name == name {
			return nil
		}
	}

	newEntry := MetadataParserEntry{Name: name, Parser: parser}

	var err error
	switch position {
	case "prepend":
		mp.prepend(newEntry)
		slog.Info("Parser registered at start", "name", name)
	case "append":
		mp.append(newEntry)
		slog.Info("Parser registered at end", "name", name)
	case "before", "after":
		err = mp.insertRelative(newEntry, position, relativeTo)
		if err != nil {
			slog.Error("Error registering parser", "error", err)
			return err
		}
		slog.Info("Parser registered", "name", name, "position", position, "relativeTo", relativeTo)
	}
	return nil
}

func (mp *MetadataParsers) prepend(entry MetadataParserEntry) {
	mp.registeredParsers = append([]MetadataParserEntry{entry}, mp.registeredParsers...)
}

func (mp *MetadataParsers) append(entry MetadataParserEntry) {
	mp.registeredParsers = append(mp.registeredParsers, entry)
}

func (mp *MetadataParsers) insertRelative(entry MetadataParserEntry, position string, relativeTo string) error {
	index := mp.findParserIndex(relativeTo)
	if index == -1 {
		return fmt.Errorf("Relative parser '%s' not found", relativeTo)
	}
	offset := 0
	if position == "after" {
		offset = 1
	}
	mp.registeredParsers = append(mp.registeredParsers[:index+offset], append([]MetadataParserEntry{entry}, mp.registeredParsers[index+offset:]...)...)
	return nil
}

func (mp *MetadataParsers) findParserIndex(name string) int {
	for i, entry := range mp.registeredParsers {
		if entry.Name == name {
			return i
		}
	}
	return -1
}

// ParseStructTags iterates over the struct tags of the given type and calls each of the registered metadata parsers.
// The metadata parsers are called with the field, field name, property, and schema ref as arguments.
// The metadata parsers can modify the property in place.
// The function is called recursively if the field type is a struct.
// The function returns if the type is not a struct or if the field is anonymous.
func (mp *MetadataParsers) ParseStructTags(t reflect.Type, schemaRef *openapi3.SchemaRef) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			fieldType := field.Type
			mp.ParseStructTags(fieldType, schemaRef)
			continue
		}

		jsonFieldName := field.Tag.Get("json")
		jsonFieldName = strings.Split(jsonFieldName, ",")[0]
		if jsonFieldName == "-" {
			jsonFieldName = field.Name
		}

		property := schemaRef.Value.Properties[jsonFieldName]
		if property == nil {
			slog.Warn("Property not found in schema", "property", jsonFieldName)
			continue
		}
		if field.Type.Kind() == reflect.Struct {
			mp.ParseStructTags(field.Type, property)
		}
		propertyCopy := *property
		propertyValue := *propertyCopy.Value

		if field.Type.Kind() == reflect.Struct {
			mp.ParseStructTags(field.Type, schemaRef.Value.Properties[jsonFieldName])
		}

		params := MetadataParserParams{
			Field:     field,
			FieldName: jsonFieldName,
			Property:  &propertyValue,
			SchemaRef: schemaRef,
		}

		for _, entry := range mp.registeredParsers {
			entry.Parser(params)
		}

		propertyCopy.Value = &propertyValue

		schemaRef.Value.Properties[jsonFieldName] = &propertyCopy
	}
}

// ExampleMetadataParser extracts the "example" tag from a struct field and assigns
// it to the Example property of the OpenAPI schema. If the field's type is an integer,
// it attempts to convert the example value to an integer, logging a warning if conversion fails.
func MetadataParserExample(params MetadataParserParams) {
	if exampleTag, ok := params.Field.Tag.Lookup("example"); ok {
		params.Property.Example = exampleTag
		if params.Property.Type.Is(openapi3.TypeInteger) {
			exNum, err := strconv.Atoi(exampleTag)
			if err != nil {
				slog.Warn("Example might be incorrect (should be integer)", "error", err)
			}
			params.Property.Example = exNum
		}
	}
}

// ValidationMetadataParser extracts the "validate" tag from a struct field and assigns
// any validation criteria to the OpenAPI schema. It supports the following validation
// criteria:
//
// - required: adds the field to the Required list on the schema
// - min=<integer>: sets the minimum value for the field in the schema
// - max=<integer>: sets the maximum value for the field in the schema
//
// The min and max values are interpreted as integers if the property type is an
// integer, and as string lengths if the property type is a string.
func MetadataParserValidation(params MetadataParserParams) {
	validateTag, ok := params.Field.Tag.Lookup("validate")
	validateTags := strings.Split(validateTag, ",")
	if ok && slices.Contains(validateTags, "required") {
		params.SchemaRef.Value.Required = append(params.SchemaRef.Value.Required, params.FieldName)
	}
	for _, validateTag := range validateTags {
		if strings.HasPrefix(validateTag, "min=") {
			min, err := strconv.Atoi(strings.Split(validateTag, "=")[1])
			if err != nil {
				slog.Warn("Min might be incorrect (should be integer)", "error", err)
			}

			if params.Property.Type.Is(openapi3.TypeInteger) {
				minPtr := float64(min)
				params.Property.Min = &minPtr
			} else if params.Property.Type.Is(openapi3.TypeString) {

				params.Property.MinLength = uint64(min)
			}
		}
		if strings.HasPrefix(validateTag, "max=") {
			max, err := strconv.Atoi(strings.Split(validateTag, "=")[1])
			if err != nil {
				slog.Warn("Max might be incorrect (should be integer)", "error", err)
			}
			if params.Property.Type.Is(openapi3.TypeInteger) {
				maxPtr := float64(max)
				params.Property.Max = &maxPtr
			} else if params.Property.Type.Is(openapi3.TypeString) {

				maxPtr := uint64(max)
				params.Property.MaxLength = &maxPtr
			}
		}
	}
}

// DescriptionMetadataParser extracts the "description" tag from a struct field and assigns
// it to the Description property of the OpenAPI schema.
func MetadataParserDescription(params MetadataParserParams) {
	description, ok := params.Field.Tag.Lookup("description")
	if ok {
		params.Property.Description = description
	}
}

// XMLMetadataParser extracts the "xml" tag from a struct field and sets the XML
// representation of the OpenAPI schema. If the xml tag is "-", the function returns
// without making any changes. The XML field name is determined from the tag or defaults
// to the field's name. The function also checks for additional XML attributes such as
// "attr" and "wrapped" to set the corresponding properties in the OpenAPI XML schema.
func MetadataParserXML(params MetadataParserParams) {
	xmlField := params.Field.Tag.Get("xml")
	if xmlField == "-" || xmlField == "" {
		return
	}
	xmlFieldName := strings.Split(xmlField, ",")[0]

	if xmlFieldName == "" {
		xmlFieldName = params.Field.Name
	}

	params.Property.XML = &openapi3.XML{
		Name: xmlFieldName,
	}

	xmlFields := strings.Split(xmlField, ",")
	if slices.Contains(xmlFields, "attr") {
		params.Property.XML.Attribute = true
	}

	if slices.Contains(xmlFields, "wrapped") {
		params.Property.XML.Wrapped = true
	}
}

// JSONMetadataParser sets the Nullable property of the OpenAPI schema based on the "json" tag.
// If the "json" tag contains ",omitempty", the Nullable property is set to true.
func MetadataParserJSON(params MetadataParserParams) {
	jsonTag, ok := params.Field.Tag.Lookup("json")
	if ok {
		if strings.Contains(jsonTag, ",omitempty") {
			params.Property.Nullable = true
		}
	}
}
