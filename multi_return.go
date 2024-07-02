package fuego

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
)

// DataOrTemplate is a struct that can return either data or a template
// depending on the asked type.
type DataOrTemplate[T any] struct {
	Data     T
	Template any
}

var (
	_ CtxRenderer    = DataOrTemplate[any]{} // Can render HTML
	_ json.Marshaler = DataOrTemplate[any]{} // Can render JSON
	_ xml.Marshaler  = DataOrTemplate[any]{} // Can render XML
)

func (m DataOrTemplate[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Data)
}

func (m DataOrTemplate[T]) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	return e.Encode(m.Data)
}

func (m DataOrTemplate[T]) Render(c context.Context, w io.Writer) error {
	switch m.Template.(type) {
	case CtxRenderer:
		return m.Template.(CtxRenderer).Render(c, w)
	case Renderer:
		return m.Template.(Renderer).Render(w)
	default:
		panic("template must be either CtxRenderer or Renderer")
	}
}

func DataOrHTML[T any](data T, template any) DataOrTemplate[T] {
	return DataOrTemplate[T]{
		Data:     data,
		Template: template,
	}
}
