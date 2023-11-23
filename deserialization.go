package fuego

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/gorilla/schema"
)

// InTransformer is an interface for entities that can be transformed.
// Useful for example for trimming strings, changing case, etc.
// Can also raise an error if the entity is not valid.
type InTransformer interface {
	InTransform() error // InTransforms the entity.
}

var ReadOptions = readOptions{
	DisallowUnknownFields: true,
	MaxBodySize:           maxBodySize,
}

// ReadJSON reads the request body as JSON.
// Can be used independantly from Fuego framework.
// Customisable by modifying ReadOptions.
func ReadJSON[B any](input io.Reader) (B, error) {
	return readJSON[B](input, ReadOptions)
}

// readJSON reads the request body as JSON.
// Can be used independantly from framework using ReadJSON,
// or as a method of Context.
// It will also read strings.
func readJSON[B any](input io.Reader, options readOptions) (B, error) {
	var body B

	// Deserialize the request body.
	dec := json.NewDecoder(input)
	if options.DisallowUnknownFields {
		dec.DisallowUnknownFields()
	}
	err := dec.Decode(&body)
	if err != nil {
		return body, fmt.Errorf("cannot decode request body: %w", err)
	}
	slog.Debug("Decoded body", "body", body)

	body, err = transform(body)
	if err != nil {
		return body, fmt.Errorf("cannot transform request body: %w", err)
	}

	err = validate(body)
	if err != nil {
		return body, fmt.Errorf("cannot validate request body: %w", err)
	}

	return body, nil
}

// ReadString reads the request body as string.
// Can be used independantly from Fuego framework.
// Customisable by modifying ReadOptions.
func ReadString[B ~string](input io.Reader) (B, error) {
	return readString[B](input, ReadOptions)
}

func readString[B ~string](input io.Reader, options readOptions) (B, error) {
	// Read the request body.
	readBody, err := io.ReadAll(input)
	if err != nil {
		return "", fmt.Errorf("cannot read request body: %w", err)
	}

	body := B(readBody)
	slog.Debug("Read body", "body", body)

	return transform(body)
}

var decoder = schema.NewDecoder()

// ReadURLEncoded reads the request body as HTML Form.
func ReadURLEncoded[B any](r *http.Request) (B, error) {
	return readURLEncoded[B](r, ReadOptions)
}

// readURLEncoded reads the request body as HTML Form.
// Can be used independantly from framework using [ReadURLEncoded],
// or as a method of Context.
func readURLEncoded[B any](r *http.Request, options readOptions) (B, error) {
	var body B

	err := r.ParseForm()
	if err != nil {
		return body, fmt.Errorf("cannot parse form: %w", err)
	}

	decoder.IgnoreUnknownKeys(!options.DisallowUnknownFields)

	err = decoder.Decode(&body, r.PostForm)
	if err != nil {
		return body, fmt.Errorf("cannot decode request body: %w", err)
	}
	slog.Debug("Decoded body", "body", body)

	body, err = transform(body)
	if err != nil {
		return body, fmt.Errorf("cannot transform request body: %w", err)
	}

	err = validate(body)
	if err != nil {
		return body, fmt.Errorf("cannot validate request body: %w", err)
	}

	return body, nil
}

// transforms the input if possible.
func transform[B any](body B) (B, error) {
	if inTransformerBody, ok := any(&body).(InTransformer); ok {
		err := inTransformerBody.InTransform()
		if err != nil {
			return body, fmt.Errorf("cannot transform request body: %w", err)
		}
		body = *any(inTransformerBody).(*B)

		slog.Debug("InTransformd body", "body", body)
	}

	return body, nil
}
