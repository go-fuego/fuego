package fuego

import (
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/gorilla/schema"
	"gopkg.in/yaml.v3"
)

// InTransformer is an interface for entities that can be transformed.
// Useful for example for trimming strings, changing case, etc.
// Can also raise an error if the entity is not valid.
type InTransformer interface {
	InTransform(context.Context) error // InTransforms the entity.
}

var ReadOptions = readOptions{
	DisallowUnknownFields: true,
	MaxBodySize:           maxBodySize,
}

// ReadJSON reads the request body as JSON.
// Can be used independently of Fuego framework.
// Customizable by modifying ReadOptions.
func ReadJSON[B any](ctx context.Context, input io.Reader) (B, error) {
	return readJSON[B](ctx, input, ReadOptions)
}

// readJSON reads the request body as JSON.
// Can be used independently of framework using ReadJSON,
// or as a method of Context.
// It will also read strings.
func readJSON[B any](ctx context.Context, input io.Reader, options readOptions) (B, error) {
	// Deserialize the request body.
	dec := json.NewDecoder(input)
	if options.DisallowUnknownFields {
		dec.DisallowUnknownFields()
	}

	return read[B](ctx, dec)
}

// ReadXML reads the request body as XML.
// Can be used independently of Fuego framework.
// Customizable by modifying ReadOptions.
func ReadXML[B any](ctx context.Context, input io.Reader) (B, error) {
	return readXML[B](ctx, input, ReadOptions)
}

// readXML reads the request body as XML.
// Can be used independently of framework using readXML,
// or as a method of Context.
func readXML[B any](ctx context.Context, input io.Reader, options readOptions) (B, error) {
	dec := xml.NewDecoder(input)
	if options.DisallowUnknownFields {
		dec.Strict = true
	}

	return read[B](ctx, dec)
}

// ReadYAML reads the request body as YAML.
// Can be used independently of Fuego framework.
// Customizable by modifying ReadOptions.
func ReadYAML[B any](ctx context.Context, input io.Reader) (B, error) {
	return readYAML[B](ctx, input, ReadOptions)
}

// readYAML reads the request body as YAML.
// Can be used independently of framework using ReadYAML,
// or as a method of Context.
func readYAML[B any](ctx context.Context, input io.Reader, options readOptions) (B, error) {
	dec := yaml.NewDecoder(input)
	if options.DisallowUnknownFields {
		dec.KnownFields(true)
	}

	return read[B](ctx, dec)
}

type decoder interface {
	Decode(v any) error
}

func read[B any](ctx context.Context, dec decoder) (B, error) {
	var body B

	err := dec.Decode(&body)
	if err != nil && !errors.Is(err, io.EOF) {
		return body, BadRequestError{
			Title:  "Decoding Failed",
			Err:    err,
			Detail: "cannot decode request body: " + err.Error(),
		}
	}
	slog.DebugContext(ctx, "Decoded body", "body", body)

	return TransformAndValidate(ctx, body)
}

// ReadString reads the request body as string.
// Can be used independently of Fuego framework.
// Customizable by modifying ReadOptions.
func ReadString[B ~string](ctx context.Context, input io.Reader) (B, error) {
	return readString[B](ctx, input, ReadOptions)
}

func readString[B ~string](ctx context.Context, input io.Reader, _ readOptions) (B, error) {
	// Read the request body.
	readBody, err := io.ReadAll(input)
	if err != nil {
		return "", BadRequestError{
			Err:    err,
			Detail: "cannot read request body: " + err.Error(),
		}
	}

	body := B(readBody)
	slog.DebugContext(ctx, "Read body", "body", body)

	return transform(ctx, body)
}

func convertSQLNullString(value string) reflect.Value {
	v := sql.NullString{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

func convertSQLNullBool(value string) reflect.Value {
	v := sql.NullBool{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}
	return reflect.ValueOf(v)
}

func newDecoder() *schema.Decoder {
	decoder := schema.NewDecoder()
	decoder.RegisterConverter(sql.NullString{}, convertSQLNullString)
	decoder.RegisterConverter(sql.NullBool{}, convertSQLNullBool)
	return decoder
}

// ReadURLEncoded reads the request body as HTML Form.
func ReadURLEncoded[B any](r *http.Request) (B, error) {
	return readURLEncoded[B](r, ReadOptions)
}

// readURLEncoded reads the request body as HTML Form.
// Can be used independently of framework using [ReadURLEncoded],
// or as a method of Context.
func readURLEncoded[B any](r *http.Request, options readOptions) (B, error) {
	var body B

	err := r.ParseForm()
	if err != nil {
		return body, fmt.Errorf("cannot parse form: %w", err)
	}

	decoder := newDecoder()
	decoder.IgnoreUnknownKeys(!options.DisallowUnknownFields)

	err = decoder.Decode(&body, r.PostForm)
	if err != nil {
		return body, BadRequestError{
			Detail: "cannot decode x-www-form-urlencoded request body: " + err.Error(),
			Err:    err,
			Errors: []ErrorItem{
				{Name: "form", Reason: "check that the form is valid, and that the content-type is correct"},
			},
		}
	}
	slog.DebugContext(r.Context(), "Decoded body", "body", body)

	return TransformAndValidate(r.Context(), body)
}

func readFormData[B any](r *http.Request, options readOptions) (B, error) {
	var body B

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return body, fmt.Errorf("cannot parse form: %w", err)
	}

	decoder := newDecoder()
	decoder.IgnoreUnknownKeys(!options.DisallowUnknownFields)

	err = decoder.Decode(&body, r.PostForm)
	if err != nil {
		return body, BadRequestError{
			Detail: "cannot decode form-data request body: " + err.Error(),
			Err:    err,
			Errors: []ErrorItem{
				{Name: "form", Reason: "check that the form is valid, and that the content-type is correct"},
			},
		}
	}
	slog.DebugContext(r.Context(), "Decoded body", "body", body)

	return TransformAndValidate(r.Context(), body)
}

// transforms the input if possible.
func transform[B any](ctx context.Context, body B) (B, error) {
	if inTransformerBody, ok := any(&body).(InTransformer); ok {
		err := inTransformerBody.InTransform(ctx)
		if err != nil {
			return body, BadRequestError{
				Title:  "Transformation Failed",
				Err:    err,
				Detail: "cannot transform request body: " + err.Error(),
				Errors: []ErrorItem{
					{Name: "transformation", Reason: "transformation failed"},
				},
			}
		}
		body = *any(inTransformerBody).(*B)

		slog.DebugContext(ctx, "InTransformd body", "body", body)
	}

	return body, nil
}

func TransformAndValidate[B any](ctx context.Context, body B) (B, error) {
	body, err := transform(ctx, body)
	if err != nil {
		return body, err
	}

	err = validate(body)
	if err != nil {
		return body, err
	}

	return body, nil
}
