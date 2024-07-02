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

func SendError(w http.ResponseWriter, r *http.Request, err error) {
	accept := parseAcceptHeader(r.Header.Get("Accept"), nil)
	if accept == "" {
		accept = "application/json"
	}

	switch accept {
	case "application/xml":
		SendXMLError(w, err)
	case "text/html":
		_ = SendHTMLError(r.Context(), w, err)
	case "text/plain":
		_ = SendText(w, err)
	case "application/json":
		SendJSONError(w, err)
	default:
		SendJSONError(w, err)
	}
}

// SendJSONError sends a JSON error response.
// If the error implements ErrorWithStatus, the status code will be set.
func SendJSONError(w http.ResponseWriter, err error) {
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
	SendJSON(w, err)
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
		status = errorStatus.StatusCode()
	}

	w.WriteHeader(status)
	SendXML(w, err)
}

func SendHTMLError(ctx context.Context, w http.ResponseWriter, err error) error {
	status := http.StatusInternalServerError
	var errorStatus ErrorWithStatus
	if errors.As(err, &errorStatus) {
		status = errorStatus.StatusCode()
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	return SendHTML(ctx, w, err.Error())
}

// SendHTML sends a HTML response.
func SendHTML(ctx context.Context, w http.ResponseWriter, ans any) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	ctxRenderer, ok := any(ans).(CtxRenderer)
	if ok {
		return ctxRenderer.Render(ctx, w)
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

func SendText(w http.ResponseWriter, ans any) error {
	var err error
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	stringToWrite, ok := any(ans).(string)
	if !ok {
		stringToWrite = *any(ans).(*string)
	}
	_, err = w.Write([]byte(stringToWrite))

	return err
}

func InferAcceptHeaderFromType(ans any) string {
	_, ok := any(ans).(string)
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

	_, ok = any(ans).(Renderer)
	if ok {
		return "text/html"
	}

	return "application/json"
}
