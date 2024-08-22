package fuego

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
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

type Sender func(http.ResponseWriter, *http.Request, any) error

// Send sends a response.
// The format is determined by the Accept header.
func Send(w http.ResponseWriter, r *http.Request, ans any) error {
	switch parseAcceptHeader(r.Header.Get("Accept"), ans) {
	case "application/xml":
		return SendXML(w, nil, ans)
	case "text/html":
		return SendHTML(w, r, ans)
	case "text/plain":
		return SendText(w, nil, ans)
	case "application/json":
		return SendJSON(w, nil, ans)
	case "application/x-yaml", "text/yaml; charset=utf-8", "application/yaml": // https://www.rfc-editor.org/rfc/rfc9512.html
		return SendYAML(w, nil, ans)
	}
	return errors.New("unsupported Accept header")
}

// SendYAML sends a YAML response.
// Declared as a variable to be able to override it for clients that need to customize serialization.
var SendYAML = func(w http.ResponseWriter, _ *http.Request, ans any) error {
	w.Header().Set("Content-Type", "application/x-yaml")
	err := yaml.NewEncoder(w).Encode(ans)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Cannot serialize returned response to YAML", "error", err)
		_, _ = w.Write([]byte(`{"error":"Cannot serialize returned response to YAML"}`))
	}
	return err
}

// SendYAMLError sends a YAML error response.
// If the error implements ErrorWithStatus, the status code will be set.
func SendYAMLError(w http.ResponseWriter, _ *http.Request, err error) {
	status := http.StatusInternalServerError
	var errorStatus ErrorWithStatus
	if errors.As(err, &errorStatus) {
		status = errorStatus.StatusCode()
	}

	w.WriteHeader(status)
	_ = SendYAML(w, nil, err.Error())
}

// SendJSON sends a JSON response.
// Declared as a variable to be able to override it for clients that need to customize serialization.
var SendJSON = func(w http.ResponseWriter, _ *http.Request, ans any) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(ans)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Cannot serialize returned response to JSON", "error", err)
		_, _ = w.Write([]byte(`{"error":"Cannot serialize returned response to JSON"}`))
	}
	return err
}

type ErrorSender = func(http.ResponseWriter, *http.Request, error)

// SendError sends an error.
// Declared as a variable to be able to override it for clients that need to customize serialization.
var SendError = func(w http.ResponseWriter, r *http.Request, err error) {
	switch parseAcceptHeader(r.Header.Get("Accept"), nil) {
	case "application/xml":
		SendXMLError(w, nil, err)
	case "text/html":
		SendHTMLError(w, nil, err)
	case "text/plain":
		SendTextError(w, r, err)
	case "application/json":
		SendJSONError(w, nil, err)
	case "application/x-yaml", "text/yaml; charset=utf-8", "application/yaml": // https://www.rfc-editor.org/rfc/rfc9512.html
		SendYAMLError(w, nil, err)
	default:
		SendJSONError(w, r, err)
	}
}

// SendJSONError sends a JSON error response.
// If the error implements ErrorWithStatus, the status code will be set.
func SendJSONError(w http.ResponseWriter, _ *http.Request, err error) {
	status := http.StatusInternalServerError
	var errorStatus ErrorWithStatus
	if errors.As(err, &errorStatus) {
		status = errorStatus.StatusCode()
	}

	w.Header().Set("Content-Type", "application/json")

	var httpError HTTPError
	if errors.As(err, &httpError) {
		w.Header().Set("Content-Type", "application/problem+json")
	}

	w.WriteHeader(status)
	_ = SendJSON(w, nil, err)
}

// SendXML sends a XML response.
// Declared as a variable to be able to override it for clients that need to customize serialization.
var SendXML = func(w http.ResponseWriter, _ *http.Request, ans any) error {
	w.Header().Set("Content-Type", "application/xml")
	return xml.NewEncoder(w).Encode(ans)
}

// SendXMLError sends a XML error response.
// If the error implements ErrorWithStatus, the status code will be set.
func SendXMLError(w http.ResponseWriter, _ *http.Request, err error) {
	status := http.StatusInternalServerError
	var errorStatus ErrorWithStatus
	if errors.As(err, &errorStatus) {
		status = errorStatus.StatusCode()
	}

	w.WriteHeader(status)
	err = SendXML(w, nil, err)
	if err != nil {
		slog.Error("Cannot serialize returned response to XML", "error", err)
		_, _ = w.Write([]byte(`{"error":"Cannot serialize returned response to XML"}`))
	}
}

// SendHTML sends a HTML response.
// Declared as a variable to be able to override it for clients that need to customize serialization.
var SendHTML = func(w http.ResponseWriter, r *http.Request, ans any) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	ctxRenderer, ok := any(ans).(CtxRenderer)
	if ok {
		return ctxRenderer.Render(r.Context(), w)
	}

	renderer, ok := any(ans).(Renderer)
	if ok {
		return renderer.Render(w)
	}

	html, ok := any(ans).(HTML)
	if ok {
		_, err := w.Write([]byte(html))
		return err
	}

	htmlString, ok := any(ans).(string)
	if ok {
		_, err := w.Write([]byte(htmlString))
		return err
	}

	// The type cannot be converted to HTML
	return fmt.Errorf("cannot serialize HTML from type %T (not string, fuego.HTML and does not implement fuego.CtxRenderer or fuego.Renderer)", ans)
}

// SendHTMLError sends a HTML response.
// If the error implements ErrorWithStatus, the status code will be set.
func SendHTMLError(w http.ResponseWriter, _ *http.Request, err error) {
	status := http.StatusInternalServerError
	var errorStatus ErrorWithStatus
	if errors.As(err, &errorStatus) {
		status = errorStatus.StatusCode()
	}

	w.WriteHeader(status)
	_ = SendHTML(w, nil, err.Error())
}

// SendText sends a HTML response.
// Declared as a variable to be able to override it for clients that need to customize serialization.
func SendText(w http.ResponseWriter, _ *http.Request, ans any) error {
	var err error
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	stringToWrite, ok := any(ans).(string)
	if !ok {
		stringToWritePtr, okPtr := any(ans).(*string)
		if okPtr {
			stringToWrite = *stringToWritePtr
		} else {
			stringToWrite = fmt.Sprintf("%v", ans)
		}
	}
	_, err = w.Write([]byte(stringToWrite))

	return err
}

// SendTextError sends a Text response.
// If the error implements ErrorWithStatus, the status code will be set.
func SendTextError(w http.ResponseWriter, _ *http.Request, err error) {
	status := http.StatusInternalServerError
	var errorStatus ErrorWithStatus
	if errors.As(err, &errorStatus) {
		status = errorStatus.StatusCode()
	}

	w.WriteHeader(status)
	_ = SendText(w, nil, err.Error())
}

func InferAcceptHeaderFromType(ans any) string {
	_, ok := any(ans).(string)
	if ok {
		return "text/plain"
	}

	_, ok = any(ans).(*string)
	if ok {
		return "text/plain"
	}

	_, ok = any(ans).(HTML)
	if ok {
		return "text/html"
	}

	_, ok = any(ans).(CtxRenderer)
	if ok {
		return "text/html"
	}

	_, ok = any(&ans).(*CtxRenderer)
	if ok {
		return "text/html"
	}

	_, ok = any(ans).(Renderer)
	if ok {
		return "text/html"
	}

	_, ok = any(&ans).(*Renderer)
	if ok {
		return "text/html"
	}

	return "application/json"
}

func parseAcceptHeader(accept string, ans any) string {
	if strings.Index(accept, ",") > 0 {
		accept = accept[:strings.Index(accept, ",")]
	}
	if accept == "*/*" {
		accept = ""
	}
	if accept == "" {
		accept = InferAcceptHeaderFromType(ans)
	}
	return accept
}
