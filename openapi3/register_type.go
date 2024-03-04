package openapi3

import "reflect"

// Registers a type to /components/schemas and returns the schema adapted to the situation.
// usually, it will be a reference to the schema.
// Example:
//
//	type S struct {
//		A string
//	}
//	type T struct {
//		S
//		B int
//	}
//
// d.RegisterType(T{})
// // will return a schema with a reference to the schema of T
func (d *Document) RegisterType(v any) *Schema {
	if v == nil {
		return nil
	}

	kind := reflect.TypeOf(v).Kind()
	switch kind {
	case reflect.Array, reflect.Slice:
		item := reflect.New(reflect.TypeOf(v).Elem()).Interface()

		name := TagFromType(item)
		if _, ok := d.Components.Schemas[name]; !ok {
			d.Components.Schemas[name] = ToSchema(item)
		}
		return &Schema{
			Type:  "array",
			Items: &Schema{Ref: "#/components/schemas/" + name},
		}
	}

	name := TagFromType(v)
	if _, ok := d.Components.Schemas[name]; !ok {
		d.Components.Schemas[name] = ToSchema(v)
	}

	return &Schema{Ref: "#/components/schemas/" + name}
}
