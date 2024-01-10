// Description: This file contains functions for sending JSON and XML responses.
package fuego

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"log/slog"
	"net/http"
	"reflect"
)

// OutTransformer is an interface for entities that can be transformed.
// Useful for example for trimming strings, changing case, etc.
// Can also raise an error if the entity is not valid.
// Must be implemented by a POINTER RECEIVER.
// Example:
//
//	type User struct {
//		Name     string `json:"name"`
//		Password string `json:"password"`
//	}
//
// // Not (u User) but (u *User)
//
//	func (u *User) OutTransform(context.Context) error {
//		u.Name = "M. " + u.Name
//		u.Password = "*****"
//		return nil
//	}
type OutTransformer interface {
	OutTransform(context.Context) error // Transforms an entity before sending it.
}

func transformOut[T any](ctx context.Context, ans T) (T, error) {
	if reflect.TypeOf(ans).Kind() == reflect.Ptr {
		// If ans is a nil pointer, we do not want to transform it.
		if reflect.ValueOf(ans).IsNil() {
			return ans, nil
		}

		outTransformer, ok := any(ans).(OutTransformer)
		if !ok {
			return ans, nil
		}

		err := outTransformer.OutTransform(ctx)
		if err != nil {
			return ans, err
		}

		return outTransformer.(T), nil
	}

	_, ok := any(ans).(OutTransformer)
	if ok {
		return ans, errors.New("OutTransformer must be implemented by a POINTER RECEIVER. Please read the [OutTransformer] documentation")
	}

	outTransformer, ok := any(&ans).(OutTransformer)
	if !ok {
		return ans, nil
	}

	err := outTransformer.OutTransform(ctx)
	if err != nil {
		return ans, err
	}
	ans = *any(outTransformer).(*T)

	return ans, nil
}

// Send sends a string response.
func Send(w http.ResponseWriter, text string) {
	_, _ = w.Write([]byte(text))
}

// SendJSON sends a JSON response.
func SendJSON(w http.ResponseWriter, ans any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(ans)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Cannot serialize JSON", "error", err)
		_, _ = w.Write([]byte(`{"error":"Cannot serialize JSON"}`))
		return
	}
}

// SendJSONError sends a JSON error response.
// If the error implements ErrorWithStatus, the status code will be set.
func SendJSONError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	errorStatus := HTTPError{
		Message: err.Error(),
	}
	if errors.As(err, &errorStatus) {
		status = errorStatus.Status()
	}

	w.WriteHeader(status)
	SendJSON(w, errorStatus)
}

// SendXML sends a XML response.
func SendXML(w http.ResponseWriter, ans any) {
	w.Header().Set("Content-Type", "application/xml")
	err := xml.NewEncoder(w).Encode(ans)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Cannot serialize XML", "error", err)
		_, _ = w.Write([]byte(`{"error":"Cannot serialize XML"}`))
		return
	}
}

// SendXMLError sends a XML error response.
// If the error implements ErrorWithStatus, the status code will be set.
func SendXMLError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	var errorStatus ErrorWithStatus
	if errors.As(err, &errorStatus) {
		status = errorStatus.Status()
	}

	w.WriteHeader(status)
	SendXML(w, err)
}
