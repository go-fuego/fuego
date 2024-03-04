package openapi3

import "reflect"

func tagFromType(v any) string {
	if v == nil {
		return "unknown-interface"
	}

	tag := dive(reflect.TypeOf(v), 4)

	switch tag {
	case "Renderer", "CtxRenderer":
		return "HTML"
	case "NetHTTP":
		return "net/http"
	}
	return tag
}

// dive returns the name of the type of the given reflect.Type.
// If the type is a pointer, slice, array, map, channel, function, or unsafe pointer,
// it will dive into the type and return the name of the type it points to.
func dive(t reflect.Type, maxDepth int) string {
	switch t.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Map, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		if maxDepth == 0 {
			return "default"
		}
		return dive(t.Elem(), maxDepth-1)
	default:
		return t.Name()
	}
}

// Registers a type to /components/schemas and returns the schema adapted to the situation.
// usually, it will be a reference to the schema.
// Example:
// 	type S struct {
// 		A string
// 	}
// 	type T struct {
// 		S
// 		B int
// 	}

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

		name := tagFromType(item)
		if _, ok := d.Components.Schemas[name]; !ok {
			d.Components.Schemas[name] = ToSchema(item)
		}
		return &Schema{
			Type:  "array",
			Items: &Schema{Ref: "#/components/schemas/" + name},
		}
	}

	name := tagFromType(v)
	if _, ok := d.Components.Schemas[name]; !ok {
		d.Components.Schemas[name] = ToSchema(v)
	}

	return &Schema{Ref: "#/components/schemas/" + name}
}
