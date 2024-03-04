package openapi3

import "reflect"

func TagFromType(v any) string {
	if v == nil {
		return "nil"
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
