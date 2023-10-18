package op

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
)

// Normalizable is an interface for entities that can be normalized.
// Useful for example for trimming strings, add custom fields, etc.
// Can also raise an error if the entity is not valid.
type Normalizable interface {
	Normalize() error // Normalizes the entity.
}

var ReadOptions = readOptions{
	DisallowUnknownFields: true,
	MaxBodySize:           maxBodySize,
}

// ReadJSON reads the request body as JSON.
// Can be used independantly from op! framework.
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

	// Validation
	err = validate(body)
	if err != nil {
		return body, fmt.Errorf("cannot validate request body: %w", err)
	}

	// Normalize input if possible.
	if normalizableBody, ok := any(&body).(Normalizable); ok {
		err := normalizableBody.Normalize()
		if err != nil {
			return body, fmt.Errorf("cannot normalize request body: %w", err)
		}
		bodyStar, ok := any(normalizableBody).(*B)
		if !ok {
			return body, fmt.Errorf("cannot retype request body: %w",
				fmt.Errorf("normalized body is not of type %T but should be", *new(B)))
		}
		body = *bodyStar

		slog.Debug("Normalized body", "body", body)
	}

	return body, nil
}

// ReadString reads the request body as string.
// Can be used independantly from op! framework.
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

	// Normalize input if possible.
	if normalizableBody, ok := any(&body).(Normalizable); ok {
		err := normalizableBody.Normalize()
		if err != nil {
			return body, fmt.Errorf("cannot normalize request body: %w", err)
		}
		bodyStar, ok := any(normalizableBody).(*B)
		if !ok {
			return body, fmt.Errorf("cannot retype request body: %w",
				fmt.Errorf("normalized body is not of type %T but should be", *new(B)))
		}
		body = *bodyStar

		slog.Debug("Normalized body", "body", body)
	}

	return body, nil
}
