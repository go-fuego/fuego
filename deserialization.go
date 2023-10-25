package op

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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

	// InTransform input if possible.
	if inTransformerBody, ok := any(&body).(InTransformer); ok {
		err := inTransformerBody.InTransform()
		if err != nil {
			return body, fmt.Errorf("cannot transform request body: %w", err)
		}
		bodyStar, ok := any(inTransformerBody).(*B)
		if !ok {
			return body, fmt.Errorf("cannot retype request body: %w",
				fmt.Errorf("transformed body is not of type %T but should be", *new(B)))
		}
		body = *bodyStar

		slog.Debug("InTransformd body", "body", body)
	}

	// Validation
	err = validate(body)
	if err != nil {
		return body, fmt.Errorf("cannot validate request body: %w", err)
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

	// InTransform input if possible.
	if inTransformerBody, ok := any(&body).(InTransformer); ok {
		err := inTransformerBody.InTransform()
		if err != nil {
			return body, fmt.Errorf("cannot transform request body: %w", err)
		}
		bodyStar, ok := any(inTransformerBody).(*B)
		if !ok {
			return body, fmt.Errorf("cannot retype request body: %w",
				fmt.Errorf("transformd body is not of type %T but should be", *new(B)))
		}
		body = *bodyStar

		slog.Debug("InTransformd body", "body", body)
	}

	return body, nil
}
